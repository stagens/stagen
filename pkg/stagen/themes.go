package stagen

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pixality-inc/golang-core/util"
	"gopkg.in/yaml.v3"
)

var (
	ErrThemeConfigNotFound = errors.New("theme config not found")
	ErrThemeAlreadyExists  = errors.New("theme already exists")
)

// nolint:unused
func (s *Impl) themesDir() string {
	dir := s.config.ThemesDir()
	if dir == "" {
		return filepath.Join(s.workDir(), "themes")
	}

	return dir
}

func (s *Impl) loadThemes(ctx context.Context) error {
	themesNamesMap := make(map[string]struct{}, 0)

	for _, page := range s.pages {
		themeName := page.Theme()
		if themeName == "" {
			themeName = s.siteConfig.Template().Theme()
		}

		themesNamesMap[themeName] = struct{}{}
	}

	for themeName, _ := range themesNamesMap {
		if err := s.loadTheme(ctx, themeName); err != nil {
			return fmt.Errorf("%w: %s: %w", ErrLoadTheme, themeName, err)
		}
	}

	return nil
}

func (s *Impl) loadTheme(ctx context.Context, themeName string) error {
	if themeName == "" {
		return ErrNoName
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading theme '%s'...", themeName)

	themeDir := filepath.Join(s.themesDir(), themeName)

	themeConfig, err := s.getThemeConfig(ctx, themeDir)
	if err != nil {
		return fmt.Errorf("theme '%s': %w", themeName, err)
	}

	if err = s.addTheme(themeName, themeDir, themeConfig); err != nil {
		return fmt.Errorf("can't add theme '%s': %w", themeName, err)
	}

	return nil
}

func (s *Impl) getThemeConfig(ctx context.Context, themeDir string) (ThemeConfig, error) {
	configFiles := s.getPossibleConfigFilenames()

	for _, configFilename := range configFiles {
		configFilePath := filepath.Join(themeDir, configFilename)

		if _, exists := util.FileExists(configFilePath); !exists {
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

	var themeConfigYaml *ThemeConfigYaml

	if err = yaml.Unmarshal(configContent, &themeConfigYaml); err != nil {
		return nil, fmt.Errorf("failed to parse theme config file %s: %w", filename, err)
	}

	return themeConfigYaml, nil
}

func (s *Impl) addTheme(
	themeId string,
	themeDir string,
	themeConfig ThemeConfig,
) error {
	if _, ok := s.pages[themeId]; ok {
		return fmt.Errorf("%w: %s", ErrThemeAlreadyExists, themeId)
	}

	layoutsIncludePaths := make([]string, 0)
	includePaths := make([]string, 0)

	templatesDir := s.templatesDir()

	layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(templatesDir, "layouts"))
	includePaths = append(includePaths, filepath.Join(templatesDir, "includes"))

	for _, extension := range s.extensions {
		layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(extension.Path(), "layouts"))
		includePaths = append(includePaths, filepath.Join(extension.Path(), "includes"))
	}

	layoutsIncludePaths = append(layoutsIncludePaths, filepath.Join(themeDir, "layouts"))
	includePaths = append(includePaths, filepath.Join(themeDir, "includes"))

	s.themes[themeId] = NewTheme(themeId, themeConfig, layoutsIncludePaths, includePaths)

	return nil
}
