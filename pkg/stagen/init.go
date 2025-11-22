package stagen

import (
	"context"
	"fmt"
)

func (s *Impl) init(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.initialized {
		return nil
	}

	log := s.log.GetLogger(ctx)

	log.Info("Initializing stagen...")

	if err := s.loadExtensions(ctx); err != nil {
		return fmt.Errorf("%w: error loading extensions: %w", ErrInit, err)
	}

	if err := s.loadDatabases(ctx); err != nil {
		return fmt.Errorf("%w: error loading databases: %w", ErrInit, err)
	}

	if err := s.loadPages(ctx); err != nil {
		return fmt.Errorf("%w: error loading pages: %w", ErrInit, err)
	}

	if err := s.loadThemes(ctx); err != nil {
		return fmt.Errorf("%w: error loading themes: %w", ErrInit, err)
	}

	if err := s.loadAggDicts(ctx); err != nil {
		return fmt.Errorf("%w: error loading agg dicts: %w", ErrInit, err)
	}

	log.Info("Initialization complete")

	s.initialized = true

	return nil
}
