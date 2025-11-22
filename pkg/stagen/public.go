package stagen

import (
	"context"
	"path/filepath"
)

func (s *Impl) publicDir() string {
	dir := s.config.PublicDir()
	if dir == "" {
		return filepath.Join(s.workDir(), "public")
	}

	return dir
}

func (s *Impl) copyPublicFiles(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	publicDir := s.publicDir()

	// @todo remove
	if publicDir == "" {
		log.Info("public dir not found")
	}

	log.Info("Copying public files...")

	// @todo copy from themes
	// @todo copy from extensions
	// @todo copy from public
	return nil
}
