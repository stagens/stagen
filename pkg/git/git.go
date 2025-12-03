package git

import (
	"context"
	"fmt"

	"github.com/pixality-inc/golang-core/cli"
	"github.com/pixality-inc/golang-core/logger"
)

type Git interface {
	HasGit(ctx context.Context) bool
	Init(ctx context.Context, workDir string) error
	Clone(ctx context.Context, workDir string, url string) error
	SubmoduleAdd(ctx context.Context, workDir string, url string, dest string) error
}

type Impl struct {
	log logger.Loggable
	cli cli.Cli
}

func New(toolPath string) *Impl {
	log := logger.NewLoggableImplWithService("git")

	return &Impl{
		log: log,
		cli: cli.New(log, toolPath),
	}
}

func (g *Impl) HasGit(ctx context.Context) bool {
	_, err := g.exec(ctx, "", "--version")

	return err == nil
}

func (g *Impl) Init(ctx context.Context, workDir string) error {
	_, err := g.exec(ctx, workDir, "init")

	return err
}

func (g *Impl) Clone(ctx context.Context, workDir string, url string) error {
	_, err := g.exec(ctx, workDir, "clone", url)

	return err
}

func (g *Impl) SubmoduleAdd(ctx context.Context, workDir string, url string, dest string) error {
	_, err := g.exec(ctx, workDir, "submodule", "add", url, dest)

	return err
}

func (g *Impl) exec(ctx context.Context, workDir string, args ...string) (string, error) {
	result, err := g.cli.Exec(ctx, args, cli.WithWorkDir(workDir))
	if err != nil {
		stderr := ""

		if result != nil {
			stderr = string(result.Stderr())
		}

		return "", fmt.Errorf("git exec failed: %w: %s", err, stderr)
	}

	return string(result.Stdout()), nil
}
