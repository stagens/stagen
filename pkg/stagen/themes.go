package stagen

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrThemeConfigNotFound = errors.New("theme config not found")
	ErrThemeNotFound       = errors.New("theme not found")
)

func (s *Impl) themesDir() string {
	return filepath.Join(s.workDir, "themes")
}

func (s *Impl) loadTheme(ctx context.Context, themeId string) (Theme, error) {
	if themeId == "" {
		return nil, ErrNoName
	}

	if existsTheme, ok := s.themes[themeId]; ok {
		return existsTheme, nil
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading theme '%s'...", themeId)

	themeDir := filepath.Join(s.themesDir(), themeId)

	themeConfig, err := s.getThemeConfig(ctx, themeDir)
	if err != nil {
		return nil, fmt.Errorf("theme '%s': %w", themeId, err)
	}

	theme, err := s.addTheme(themeId, themeDir, themeConfig)
	if err != nil {
		return nil, fmt.Errorf("can't add theme '%s': %w", themeId, err)
	}

	return theme, nil
}

func (s *Impl) getThemeConfig(ctx context.Context, themeDir string) (ThemeConfig, error) {
	configFiles := s.getPossibleConfigFilenames()

	for _, configFilename := range configFiles {
		configFilePath := filepath.Join(themeDir, configFilename)

		if exists, err := s.storage.FileExists(ctx, configFilePath); err != nil {
			return nil, fmt.Errorf("faile to check if file %s exists: %w", configFilePath, err)
		} else if !exists {
			continue
		}

		themeConfig, err := s.readThemeConfig(ctx, configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read theme config %s: %w", configFilePath, err)
		}

		return themeConfig, nil
	}

	return nil, ErrThemeConfigNotFound
}

func (s *Impl) readThemeConfig(ctx context.Context, filename string) (ThemeConfig, error) {
	configContent, err := s.readFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir config file %s: %w", filename, err)
	}

	var themeConfigYaml ThemeConfigYaml

	if err = yaml.Unmarshal(configContent, &themeConfigYaml); err != nil {
		return nil, fmt.Errorf("failed to parse theme config file %s: %w", filename, err)
	}

	return &themeConfigYaml, nil
}

func (s *Impl) addTheme(
	themeId string,
	themeDir string,
	themeConfig ThemeConfig,
) (Theme, error) {
	layoutsIncludePaths := make([]string, 0)
	importPaths := make([]string, 0)
	includePaths := make([]string, 0)

	templatesDir := s.templatesDir()

	layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(templatesDir, "layouts"))
	importPaths = append(importPaths, filepath.Join(templatesDir, "imports"))
	includePaths = append(includePaths, filepath.Join(templatesDir, "includes"))

	for _, extension := range s.extensions {
		layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(extension.Path(), "layouts"))
		importPaths = append(importPaths, filepath.Join(extension.Path(), "imports"))
		includePaths = append(includePaths, filepath.Join(extension.Path(), "includes"))
	}

	layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(themeDir, "layouts"))
	importPaths = append(importPaths, filepath.Join(themeDir, "imports"))
	includePaths = append(includePaths, filepath.Join(themeDir, "includes"))

	s.themes[themeId] = NewTheme(themeId, themeDir, themeConfig, s.storage, layoutsIncludePaths, importPaths, includePaths)

	return s.themes[themeId], nil
}
