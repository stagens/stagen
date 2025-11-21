package stagen

import (
	"context"
	"path/filepath"
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
	if false {
		if err := s.loadTheme(ctx, "asdasd"); err != nil {
			return err
		}
	}

	// @todo
	return nil
}

func (s *Impl) loadTheme(ctx context.Context, themeName string) error {
	if themeName == "" {
		return ErrNoName
	}

	log := s.log.GetLogger(ctx)

	log.Infof("Loading theme '%s'...", themeName)

	// @todo
	return nil
}
