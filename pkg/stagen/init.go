package stagen

import (
	"context"
	"fmt"
	"os"
	"time"
)

func (s *Impl) init(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.initialized {
		return nil
	}

	log := s.log.GetLogger(ctx)

	log.Info("Initializing stagen...")

	log.Infof("Creating build dir...")

	buildDir := s.buildDir()

	if err := os.MkdirAll(buildDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create build dir: %w", err)
	}

	s.buildTime = time.Now()

	if err := s.loadExtensions(ctx); err != nil {
		return fmt.Errorf("%w: error loading extensions: %w", ErrInit, err)
	}

	if err := s.loadDatabases(ctx); err != nil {
		return fmt.Errorf("%w: error loading databases: %w", ErrInit, err)
	}

	if err := s.loadPages(ctx); err != nil {
		return fmt.Errorf("%w: error loading pages: %w", ErrInit, err)
	}

	if err := s.loadAggDicts(ctx); err != nil {
		return fmt.Errorf("%w: error loading agg dicts: %w", ErrInit, err)
	}

	log.Info("Initialization complete")

	s.initialized = true

	return nil
}
