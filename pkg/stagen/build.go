package stagen

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/pixality-inc/golang-core/timetrack"
	"github.com/pixality-inc/golang-core/util"
)

type PageRenderConfig struct {
	Page    Page
	Theme   Theme
	Data    map[string]any
	Content []byte
}

func (s *Impl) Build(ctx context.Context) error {
	if err := s.init(ctx); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	log := s.log.GetLogger(ctx)

	log.Info("Running build...")

	track := timetrack.New(ctx)

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
		"base",
		templateConfig.Theme(),
		templateConfig.DefaultLayout(),
		"",
		false,
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

func (s *Impl) buildDir() string {
	return filepath.Join(s.workDir, "build")
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

	pageConfig := page.Config()

	log.Infof(
		"Building page '%s' (theme: %s, layout: %s, isHidden: %v, isDraft: %v)...",
		pageId,
		pageRenderConfig.Theme.Name(),
		pageConfig.Layout(),
		pageConfig.IsHidden(),
		pageConfig.IsDraft(),
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

	data, err := s.getTemplateData(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to get template data for page '%s': %w", pageId, err)
	}

	pageContent := page.Content()

	pageRenderConfig := &PageRenderConfig{
		Page:    page,
		Theme:   theme,
		Data:    data,
		Content: pageContent,
	}

	return pageRenderConfig, nil
}

func (s *Impl) renderPage(ctx context.Context, pageRenderConfig *PageRenderConfig) ([]byte, error) {
	pageId := pageRenderConfig.Page.Id()
	pageConfig := pageRenderConfig.Page.Config()

	renderedContent, err := pageRenderConfig.Theme.Render(
		ctx,
		pageConfig.Imports(),
		pageConfig.Layout(),
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
) (map[string]any, error) {
	pageData := func(pageEntry Page) (map[string]any, error) {
		pageId := pageEntry.Id()
		pageConfig := pageEntry.Config()
		pageFileInfo := pageEntry.FileInfo()

		pageUrl, err := url.JoinPath(s.siteConfig.BaseUrl(), pageEntry.Uri())
		if err != nil {
			return nil, fmt.Errorf("failed to resolve page '%s' url: %w", pageId, err)
		}

		data := map[string]any{
			"Id":         pageId,
			"Name":       pageEntry.Name(),
			"Uri":        pageEntry.Uri(),
			"Url":        pageUrl,
			"Title":      pageConfig.Title(),
			"IsHidden":   pageConfig.IsHidden(),
			"IsDraft":    pageConfig.IsDraft(),
			"IsSystem":   pageConfig.IsSystem(),
			"CreatedAt":  pageFileInfo.CreatedAt,
			"ModifiedAt": pageFileInfo.ModifiedAt,
			"AccessedAt": pageFileInfo.AccessedAt,
			"ChangedAt":  pageFileInfo.ChangedAt,
			"Variables":  pageConfig.Variables(),
			"Imports":    pageConfig.Imports(),
			"Includes":   pageConfig.Includes(),
			"Extras":     pageConfig.Extras(),
		}

		return data, nil
	}

	pagesData := make(map[string]any)

	for _, pageEntry := range s.pages {
		pageEntryData, err := pageData(pageEntry)
		if err != nil {
			return nil, fmt.Errorf("failed to get page entry data for page '%s': %w", pageEntry.Id(), err)
		}

		pagesData[pageEntry.Id()] = pageEntryData
	}

	pageEntryData, err := pageData(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get page entry data for page '%s': %w", page.Id(), err)
	}

	data := map[string]any{
		"Site": map[string]any{
			"Name":      s.siteConfig.Name(),
			"BaseUrl":   s.siteConfig.BaseUrl(),
			"Author":    s.siteConfig.Author(),
			"Copyright": s.siteConfig.Copyright(),
			"Logo":      s.siteConfig.Logo(),
		},
		"Page": pageEntryData,
		"System": map[string]any{
			"BuildTime": s.buildTime,
			"Now":       time.Now(),
		},
		"Pages":        pagesData,
		"Databases":    s.databases,
		"AggDicts":     s.aggDicts,
		"AggDictsData": s.aggDictsData,
	}

	for k, v := range page.Config().Variables() {
		data[k] = v //nolint:modernize // @todo
	}

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

		if err := s.storage.MkDir(ctx, dir); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}

		s.createdDirs[dir] = struct{}{}
	}

	if err := s.storage.Write(ctx, saveFilename, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("failed to save file '%s': %w", saveFilename, err)
	}

	return nil
}
