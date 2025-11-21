package stagen

import (
	"context"
	"fmt"
	"path/filepath"
)

func (s *Impl) Build(ctx context.Context) error {
	if err := s.init(ctx); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if err := s.build(ctx); err != nil {
		return fmt.Errorf("failed to build: %w", err)
	}

	if err := s.copyPublicFiles(ctx); err != nil {
		return fmt.Errorf("failed to copy public files: %w", err)
	}

	return nil
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

	return nil
}

func (s *Impl) copyPublicFiles(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Copying public files...")

	return nil
}
