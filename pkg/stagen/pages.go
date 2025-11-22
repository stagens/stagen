package stagen

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/pixality-inc/golang-core/util"
	"gopkg.in/yaml.v3"

	"stagen/pkg/filetree"
)

func (s *Impl) pagesDir() string {
	dir := s.config.PagesDir()
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
	dirConfigs []DirConfig,
) error {
	if pageFilename == "" {
		return ErrNoName
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading page %s...", pageFilename)

	pageDir := filepath.Dir(pageFilename)
	pageDirWithoutWorkDir := strings.TrimPrefix(pageDir, s.workDir())
	pageFilenameWithoutExt, fullExt := removeFileExtension(pageFilename)
	fileExt := strings.TrimPrefix(pageFilename, pageFilenameWithoutExt)
	templateExt := filepath.Ext(fullExt)

	if !strings.HasPrefix(fileExt, templateExt) {
		fileExt = strings.TrimSuffix(fileExt, templateExt)
	} else {
		templateExt = ""
	}

	isTemplate := templateExtensionRegexp.MatchString(templateExt)

	pageFileInfo := PageFileInfo{
		Filename:                 pageFilename,
		BaseFilename:             filepath.Base(pageFilename),
		Path:                     pageDir,
		PathWithoutWorkDir:       pageDirWithoutWorkDir,
		FilenameWithoutExtension: pageFilenameWithoutExt,
		FullExtension:            fullExt,
		FileExtension:            fileExt,
		TemplateExtension:        templateExt,
		IsTemplate:               isTemplate,
	}

	fileContent, err := s.readFile(ctx, pageFilename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var pageConfigYamlNode *yaml.Node

	fileContent, err = frontmatter.Parse(bytes.NewReader(fileContent), &pageConfigYamlNode)
	if err != nil {
		return fmt.Errorf("failed to parse file front matter: %w", err)
	}

	var (
		pageVariables  map[string]any
		pageConfigYaml *PageConfigYaml
		pageConfig     PageConfig
	)

	if pageConfigYamlNode != nil {
		if err = pageConfigYamlNode.Decode(&pageVariables); err != nil {
			return fmt.Errorf("failed to parse page variables yaml config: %w", err)
		}

		if err = pageConfigYamlNode.Decode(&pageConfigYaml); err != nil {
			return fmt.Errorf("failed to parse page yaml config: %w", err)
		}
	}

	if pageVariables == nil {
		pageVariables = make(map[string]any)
	}

	if pageConfigYaml == nil {
		pageConfig = NewPageConfigImpl()
	} else {
		pageConfig = pageConfigYaml
	}

	page := NewPage(
		pageFileInfo.PathWithoutWorkDir,
		&pageFileInfo,
		fileContent,
		pageVariables,
		dirConfigs,
		pageConfig,
	)

	s.pages = append(s.pages, page)

	return nil
}

func (s *Impl) processPagesDirEntry(ctx context.Context, dirEntry filetree.Entry, dirConfigs []DirConfig) error {
	dir := filepath.Join(dirEntry.Path(), dirEntry.Name())

	entryDirConfigs, err := s.readDirConfigs(ctx, dir)
	if err != nil {
		return fmt.Errorf("failed to read dir config: %w", err)
	}

	childDirConfigs := make([]DirConfig, 0, len(dirConfigs)+len(entryDirConfigs))

	childDirConfigs = append(childDirConfigs, entryDirConfigs...)

	if err = s.processPagesDirEntries(ctx, dirEntry.Children(), childDirConfigs); err != nil {
		return err
	}

	return nil
}

func (s *Impl) processPagesDirEntries(ctx context.Context, entries []filetree.Entry, dirConfigs []DirConfig) error {
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

func (s *Impl) readDirConfigs(ctx context.Context, dir string) ([]DirConfig, error) {
	dirConfigs := make([]DirConfig, 0)

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

		dirConfigs = append(dirConfigs, dirConfig)
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
