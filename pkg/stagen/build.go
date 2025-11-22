package stagen

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pixality-inc/golang-core/timetrack"
	"github.com/pixality-inc/golang-core/util"
)

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

	// @todo
	return nil
}

func (s *Impl) copyPublicFiles(ctx context.Context) error {
	log := s.log.GetLogger(ctx)

	log.Info("Copying public files...")

	// @todo
	return nil
}
