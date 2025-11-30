package stagen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/djherbis/times"
	"github.com/pixality-inc/golang-core/util"
	"gopkg.in/yaml.v3"

	"stagen/pkg/filetree"
)

const (
	htmlExtension     = ".html"
	markdownExtension = ".md"
)

var ErrPageAlreadyExists = errors.New("page already exists")

func (s *Impl) pagesDir() string {
	dir := s.config.Dirs().Pages()
	if dir == "" {
		return filepath.Join(s.workDir(), "pages")
	}

	return dir
}

func (s *Impl) loadPages(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Loading pages...")

	pagesDir := s.pagesDir()

	tree, err := filetree.Tree(ctx, pagesDir, filetree.NoMaxLevel)
	if err != nil {
		return fmt.Errorf("failed to build tree for pages dir: %w", err)
	}

	if err = s.processPagesDirEntry(ctx, tree, nil); err != nil {
		return fmt.Errorf("failed to process pages dir entries: %w", err)
	}

	return nil
}

func (s *Impl) loadPage(
	ctx context.Context,
	pageFilename string,
	dirConfigs []PageConfig,
) error {
	if pageFilename == "" {
		return ErrNoName
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading page '%s'...", pageFilename)

	pageFileInfo, err := s.getPageFileInfo(pageFilename)
	if err != nil {
		return fmt.Errorf("failed to get page '%s' file info: %w", pageFilename, err)
	}

	fileContent, err := s.readFile(ctx, pageFilename)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", pageFilename, err)
	}

	page, err := s.createPage(
		ctx,
		pageFileInfo,
		fileContent,
		nil,
		dirConfigs,
	)
	if err != nil {
		return fmt.Errorf("failed to create page '%s': %w", pageFilename, err)
	}

	if err = s.addPage(page); err != nil {
		return fmt.Errorf("failed to add page '%s': %w", pageFilename, err)
	}

	return nil
}

func (s *Impl) createPage(
	ctx context.Context,
	pageFileInfo *PageFileInfo,
	content []byte,
	extraVariables map[string]any,
	dirConfigs []PageConfig,
) (Page, error) {
	var (
		err            error
		pageVariables  map[string]any
		pageConfigYaml *PageConfigYaml
		readPageConfig PageConfig
	)

	pageContent := make([]byte, len(content))
	copy(pageContent, content)

	content, err = frontmatter.Parse(bytes.NewReader(pageContent), &pageVariables)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file front matter: %w", err)
	}

	_, err = frontmatter.Parse(bytes.NewReader(pageContent), &pageConfigYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file front matter: %w", err)
	}

	if pageVariables == nil {
		pageVariables = make(map[string]any)
	}

	for k, v := range extraVariables {
		pageVariables[k] = v //nolint:modernize // @todo
	}

	if pageConfigYaml == nil {
		readPageConfig = NewDefaultPageConfig("empty", pageVariables)
	} else {
		readPageConfig = pageConfigYaml.ToPageConfig(pageVariables)
	}

	basePageConfig := s.getBasePageConfig()

	tempPageConfig := MergePageConfigs(basePageConfig, readPageConfig)

	themeId := tempPageConfig.Theme()

	theme, err := s.loadTheme(ctx, themeId)
	if err != nil {
		return nil, fmt.Errorf("failed to load theme '%s': %w", themeId, err)
	}

	pageConfig := MergePageConfigs(basePageConfig, theme.Config().ToPageConfig())

	for _, dirConfig := range dirConfigs {
		pageConfig = MergePageConfigs(pageConfig, dirConfig)
	}

	pageConfig = MergePageConfigs(pageConfig, readPageConfig)

	pageName := filepath.Join(pageFileInfo.PathWithoutWorkDirAndPagesDir, pageFileInfo.FilenameWithoutExtension)

	pageId := pageName
	if pageId == "index" {
		pageId = ""
	}

	pageUri := "/" + pageId

	isIndex := strings.HasSuffix(pageUri, "/index")

	if s.config.Settings().UseUriHtmlFileExtension() {
		switch {
		case isIndex:
			pageUri += htmlExtension
		default:
			switch pageFileInfo.FileExtension {
			case htmlExtension, markdownExtension:
				pageUri += htmlExtension
			default:
				pageUri += pageFileInfo.FileExtension
			}
		}
	} else {
		switch pageFileInfo.FileExtension {
		case htmlExtension, markdownExtension:
		default:
			pageUri += pageFileInfo.FileExtension
		}
	}

	pageUri, _ = strings.CutSuffix(pageUri, "/index.html")
	pageUri, _ = strings.CutSuffix(pageUri, "/index")
	pageUri, _ = strings.CutSuffix(pageUri, "/.html")

	if pageUri == "" {
		pageUri = "/"
	}

	page := NewPage(
		pageId,
		pageName,
		pageUri,
		pageFileInfo,
		content,
		dirConfigs,
		pageConfig,
	)

	return page, nil
}

func (s *Impl) addPage(page Page) error {
	pageId := page.Id()

	if _, ok := s.pages[pageId]; ok {
		return fmt.Errorf("%w: %s", ErrPageAlreadyExists, pageId)
	}

	s.pages[pageId] = page

	return nil
}

func (s *Impl) getPageFileInfo(pageFilename string) (*PageFileInfo, error) {
	stat, err := times.Stat(pageFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to stat page file '%s': %w", pageFilename, err)
	}

	pageFileInfo := NewPageFileInfo(
		pageFilename,
		s.workDir(),
		s.pagesDir(),
		stat,
	)

	return pageFileInfo, nil
}

func (s *Impl) processPagesDirEntry(
	ctx context.Context,
	dirEntry filetree.Entry,
	dirConfigs []PageConfig,
) error {
	dir := filepath.Join(dirEntry.Path(), dirEntry.Name())

	entryDirConfigs, err := s.readDirConfigs(ctx, dir)
	if err != nil {
		return fmt.Errorf("failed to read dir config: %w", err)
	}

	childDirConfigs := make([]PageConfig, 0, len(dirConfigs)+len(entryDirConfigs))

	childDirConfigs = append(childDirConfigs, entryDirConfigs...)

	if err = s.processPagesDirEntries(ctx, dirEntry.Children(), childDirConfigs); err != nil {
		return err
	}

	return nil
}

func (s *Impl) processPagesDirEntries(
	ctx context.Context,
	entries []filetree.Entry,
	dirConfigs []PageConfig,
) error {
	log := s.log.GetLogger(ctx)

	for _, dirEntry := range entries {
		if dirEntry.IsDir() {
			if err := s.processPagesDirEntry(ctx, dirEntry, dirConfigs); err != nil {
				return err
			}

			continue
		}

		pageFilename := filepath.Join(dirEntry.Path(), dirEntry.Name())

		if pageIgnoreFilenameRegexp.MatchString(dirEntry.Name()) {
			log.Debugf("Skipping page file %s...", pageFilename)

			continue
		}

		if err := s.loadPage(ctx, pageFilename, dirConfigs); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadPage, pageFilename, err)
		}
	}

	return nil
}

func (s *Impl) readDirConfigs(ctx context.Context, dir string) ([]PageConfig, error) {
	dirConfigs := make([]PageConfig, 0)

	configFiles := s.getPossibleConfigFilenames()

	for _, configFilename := range configFiles {
		configFilePath := filepath.Join(dir, configFilename)

		if _, exists := util.FileExists(configFilePath); !exists {
			continue
		}

		dirConfig, err := s.readDirConfig(ctx, configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read dir %s config %s: %w", dir, configFilePath, err)
		}

		dirConfigs = append(dirConfigs, dirConfig.ToPageConfig(dir))
	}

	return dirConfigs, nil
}

func (s *Impl) readDirConfig(ctx context.Context, filename string) (DirConfig, error) {
	configContent, err := s.readFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir config file %s: %w", filename, err)
	}

	var dirConfigYaml *DirConfigYaml

	if err = yaml.Unmarshal(configContent, &dirConfigYaml); err != nil {
		return nil, fmt.Errorf("failed to parse dir config file %s: %w", filename, err)
	}

	return dirConfigYaml, nil
}
