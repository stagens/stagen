package cli

import (
	"context"
	"path/filepath"

	"github.com/pixality-inc/golang-core/logger"

	"stagen/internal/config"
	"stagen/pkg/stagen"
)

type Cli interface {
	Init(ctx context.Context, workDir string, name string) error
	Build(ctx context.Context, workDir string) error
	Web(ctx context.Context, workDir string) error
}

type Impl struct {
	log    logger.Loggable
	stagen stagen.Stagen
}

func New(stagenTool stagen.Stagen) *Impl {
	return &Impl{
		log:    logger.NewLoggableImplWithService("cli"),
		stagen: stagenTool,
	}
}

func (c *Impl) Init(ctx context.Context, workDir string, name string) error {
	defaultConfig := config.NewConfig(workDir)

	if err := c.init(ctx, workDir, defaultConfig); err != nil {
		return err
	}

	if err := c.stagen.NewProject(ctx, name); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Build(ctx context.Context, workDir string) error {
	if err := c.init(ctx, workDir, nil); err != nil {
		return err
	}

	if err := c.stagen.Build(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Impl) Web(ctx context.Context, workDir string) error {
	if err := c.init(ctx, workDir, nil); err != nil {
		return err
	}

	if err := c.stagen.Web(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Impl) init(ctx context.Context, workDir string, cfg *config.Config) error {
	if cfg == nil {
		var err error

		configFilename := filepath.Join(workDir, "config.yaml")

		cfg, err = config.NewConfigFromFile(configFilename)
		if err != nil {
			return err
		}
	}

	if err := c.stagen.Init(ctx, &cfg.Stagen, &cfg.Site); err != nil {
		return err
	}

	return nil
}
