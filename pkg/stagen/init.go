package stagen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"stagen/internal/build"
)

var ErrNotInitialized = errors.New("not initialized")

func (s *Impl) Init(_ context.Context, cfg Config, siteConfig SiteConfig) error {
	s.config = cfg
	s.siteConfig = siteConfig

	return nil
}

func (s *Impl) versionInfo() string {
	parts := make([]string, 0)

	parts = append(parts, "version: "+Version)

	if build.CiPipelineId != build.DefaultValue {
		parts = append(parts, "ci pipeline dd: "+build.CiPipelineId)
	}

	if build.GitTag != build.DefaultValue {
		parts = append(parts, "git tag: "+build.GitTag)
	}

	if build.GitBranch != build.DefaultValue {
		parts = append(parts, "git branch: "+build.GitBranch)
	}

	if build.GitCommit != build.DefaultValue {
		parts = append(parts, "git commit: "+build.GitCommit)
	}

	if build.GitCommitShort != build.DefaultValue {
		parts = append(parts, "git commit short: "+build.GitCommitShort)
	}

	return strings.Join(parts, ", ")
}

func (s *Impl) init(ctx context.Context) error {
	s.initMutex.Lock()
	defer s.initMutex.Unlock()

	if s.initialized {
		return nil
	}

	log := s.log.GetLogger(ctx)

	log.Info("[STAGEN] " + s.versionInfo())

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

	for range 2 {
		s.aggDictsData = make(map[string]map[string]map[string][]Page)

		if err := s.loadAggDictsData(ctx); err != nil {
			return fmt.Errorf("%w: error loading agg dicts data: %w", ErrInit, err)
		}

		s.generators = make(map[string]Generator)

		if err := s.loadGenerators(ctx); err != nil {
			return fmt.Errorf("%w: error loading generators: %w", ErrInit, err)
		}

		if err := s.buildGenerators(ctx); err != nil {
			return fmt.Errorf("failed to build generators: %w", err)
		}
	}

	log.Info("Initialization complete")

	s.initialized = true

	return nil
}
