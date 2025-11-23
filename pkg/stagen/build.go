package stagen

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/pixality-inc/golang-core/timetrack"
	"github.com/pixality-inc/golang-core/util"
)

type PageRenderConfig struct {
	Page       Page
	PageConfig PageConfig
	Theme      Theme
	Layout     string
	Data       map[string]any
	Content    []byte
}

func (s *Impl) Build(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Running build...")

	track := timetrack.New()

	if err := s.init(ctx); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if err := s.build(ctx); err != nil {
		return fmt.Errorf("failed to build: %w", err)
	}

	if err := s.copyPublicFiles(ctx); err != nil {
		return fmt.Errorf("failed to copy public files: %w", err)
	}

	log.Infof("Build finished in %s", util.FormatDuration(track.Finish()))

	return nil
}

func (s *Impl) getBasePageConfig() PageConfig {
	templateConfig := s.siteConfig.Template()

	var pageConfig PageConfig = NewPageConfig(
		templateConfig.Theme(),
		templateConfig.DefaultLayout(),
		"",
		false,
		false,
		templateConfig.Variables(),
		templateConfig.Imports(),
		templateConfig.Includes(),
		templateConfig.Extras(),
	)

	for _, extension := range s.extensions {
		pageConfig = MergePageConfigs(pageConfig, extension.Config().ToPageConfig())
	}

	return pageConfig
}

// nolint:unused
func (s *Impl) buildDir() string {
	dir := s.config.BuildDir()
	if dir == "" {
		return filepath.Join(s.workDir(), "build")
	}

	return dir
}

func (s *Impl) build(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Building pages...")

	for _, page := range s.pages {
		if err := s.buildPage(ctx, page); err != nil {
			return fmt.Errorf("failed to build page '%s': %w", page.Id(), err)
		}
	}

	return nil
}

func (s *Impl) buildPage(ctx context.Context, page Page) error {
	log := s.log.GetLogger(ctx)

	pageId := page.Name()

	pageRenderConfig, err := s.getPageRenderConfig(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to get page render config for page '%s': %w", pageId, err)
	}

	log.Infof(
		"Building page '%s' (theme: %s, layout: %s, isHidden: %v, isDraft: %v)...",
		pageId,
		pageRenderConfig.Theme.Name(),
		pageRenderConfig.Layout,
		pageRenderConfig.PageConfig.IsHidden(),
		pageRenderConfig.PageConfig.IsDraft(),
	)

	renderedContent, err := s.renderPage(ctx, pageRenderConfig)
	if err != nil {
		return fmt.Errorf("failed to render page '%s': %w", pageId, err)
	}

	pageFileInfo := page.FileInfo()

	if err = s.saveBuildPage(ctx, pageFileInfo, renderedContent); err != nil {
		return fmt.Errorf("failed to save page '%s': %w", pageId, err)
	}

	return nil
}

func (s *Impl) getPageRenderConfig(
	ctx context.Context,
	page Page,
) (*PageRenderConfig, error) {
	pageId := page.Id()

	pageConfig := page.Config()

	themeId := pageConfig.Theme()

	theme, ok := s.themes[themeId]
	if !ok {
		return nil, fmt.Errorf("%w: %s for page '%s'", ErrThemeNotFound, themeId, pageId)
	}

	data, err := s.getTemplateData(ctx, page, pageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get template data for page '%s': %w", pageId, err)
	}

	pageContent := page.Content()

	pageRenderConfig := &PageRenderConfig{
		Page:       page,
		PageConfig: pageConfig,
		Theme:      theme,
		Layout:     pageConfig.Layout(),
		Data:       data,
		Content:    pageContent,
	}

	return pageRenderConfig, nil
}

func (s *Impl) renderPage(ctx context.Context, pageRenderConfig *PageRenderConfig) ([]byte, error) {
	pageId := pageRenderConfig.Page.Id()

	renderedContent, err := pageRenderConfig.Theme.Render(
		ctx,
		pageRenderConfig.PageConfig.Imports(),
		pageRenderConfig.Layout,
		pageRenderConfig.Content,
		pageRenderConfig.Page.FileInfo().IsMarkdown,
		pageRenderConfig.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render page '%s': %w", pageId, err)
	}

	return renderedContent, nil
}

func (s *Impl) getTemplateData(
	_ context.Context,
	page Page,
	pageConfig PageConfig,
) (map[string]any, error) {
	pageId := page.Id()

	pageUrl, err := url.JoinPath(s.siteConfig.BaseUrl(), page.Uri())
	if err != nil {
		return nil, fmt.Errorf("failed to resolve page '%s' url: %w", pageId, err)
	}

	pageFileInfo := page.FileInfo()

	data := map[string]any{
		"Site": map[string]any{
			"Name":      s.siteConfig.Name(),
			"BaseUrl":   s.siteConfig.BaseUrl(),
			"Author":    s.siteConfig.Author(),
			"Copyright": s.siteConfig.Copyright(),
			"Logo":      s.siteConfig.Logo(),
		},
		"Page": map[string]any{
			"Id":         pageId,
			"Name":       page.Name(),
			"Uri":        page.Uri(),
			"Url":        pageUrl,
			"Title":      pageConfig.Title(),
			"CreatedAt":  pageFileInfo.CreatedAt,
			"ModifiedAt": pageFileInfo.ModifiedAt,
			"AccessedAt": pageFileInfo.AccessedAt,
			"ChangedAt":  pageFileInfo.ChangedAt,
			"Variables":  pageConfig.Variables(),
		},
		"System": map[string]any{
			"BuildTime": s.buildTime,
			"Now":       time.Now(),
		},
		"Pages":     s.pages,
		"Databases": s.databases,
		"AggDicts":  nil, // @todo AggDicts
	}

	maps.Copy(data, pageConfig.Variables())

	return data, nil
}

func (s *Impl) saveBuildPage(
	ctx context.Context,
	pageFileInfo *PageFileInfo,
	content []byte,
) error {
	log := s.log.GetLogger(ctx)

	fileExt := pageFileInfo.FileExtension
	if pageFileInfo.IsMarkdown {
		fileExt = ".html"
	}

	filename := filepath.Join(pageFileInfo.PathWithoutWorkDirAndPagesDir, pageFileInfo.FilenameWithoutExtension) + fileExt
	saveFilename := filepath.Join(s.buildDir(), filename)

	log.Debugf(
		"Saving %d bytes to %s...",
		len(content),
		saveFilename,
	)

	dir := filepath.Dir(saveFilename)
	if _, ok := s.createdDirs[dir]; !ok {
		log.Debugf("Creating directory '%s'...", dir)

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}

		s.createdDirs[dir] = struct{}{}
	}

	//nolint:gosec
	if err := os.WriteFile(saveFilename, content, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file '%s': %w", saveFilename, err)
	}

	return nil
}
