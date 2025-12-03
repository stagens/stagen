package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/pixality-inc/golang-core/clock"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/storage"
	"github.com/pixality-inc/golang-core/storage/providers"

	"stagen/internal/config"
	"stagen/pkg/git"
	"stagen/pkg/stagen"
)

type Cli interface {
	Init(ctx context.Context, workDir string, name string, withGit bool) error
	Build(ctx context.Context, workDir string) error
	Watch(ctx context.Context, workDir string) error
	Web(ctx context.Context, workDir string) error
	Dev(ctx context.Context, workDir string) error
}

type Impl struct {
	log   logger.Loggable
	clock clock.Clock
	git   git.Git
}

func New(clocks clock.Clock, gitTool git.Git) *Impl {
	return &Impl{
		log:   logger.NewLoggableImplWithService("cli"),
		clock: clocks,
		git:   gitTool,
	}
}

func (c *Impl) Init(ctx context.Context, workDir string, name string, withGit bool) error {
	defaultConfig := config.NewConfig()

	stagenTool, err := c.init(ctx, workDir, defaultConfig)
	if err != nil {
		return err
	}

	if err = stagenTool.NewProject(ctx, name, withGit); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Build(ctx context.Context, workDir string) error {
	stagenTool, err := c.init(ctx, workDir, nil)
	if err != nil {
		return err
	}

	if err = stagenTool.Build(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Watch(ctx context.Context, workDir string) error {
	stagenTool, err := c.init(ctx, workDir, nil)
	if err != nil {
		return err
	}

	if err = stagenTool.Watch(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Web(ctx context.Context, workDir string) error {
	stagenTool, err := c.init(ctx, workDir, nil)
	if err != nil {
		return err
	}

	if err = stagenTool.Web(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Dev(ctx context.Context, workDir string) error {
	stagenTool, err := c.init(ctx, workDir, nil)
	if err != nil {
		return err
	}

	log := c.log.GetLogger(ctx)

	if err = stagenTool.Build(ctx); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	wg := sync.WaitGroup{}

	wg.Go(func() {
		if err := stagenTool.Watch(ctx); err != nil {
			log.WithError(err).Error("Watch failed")
		}
	})

	wg.Go(func() {
		if err := stagenTool.Web(ctx); err != nil {
			log.WithError(err).Error("Web failed")
		}
	})

	wg.Wait()

	return nil
}

func (c *Impl) init(_ context.Context, workDir string, cfg *config.Config) (stagen.Stagen, error) {
	if cfg == nil {
		var err error

		configFilename := filepath.Join(workDir, "config.yaml")

		cfg, err = config.NewConfigFromFile(configFilename)
		if err != nil {
			return nil, err
		}
	}

	localStorage := storage.NewLocalStorage(
		providers.NewOsProvider(workDir),
		providers.NoUrlProviderImpl,
	)

	stagenTool := stagen.New(&cfg.Stagen, &cfg.Site, c.clock, c.git, localStorage, workDir)

	return stagenTool, nil
}
