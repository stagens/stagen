package wiring

import (
	"context"
	"time"

	"github.com/pixality-inc/golang-core/base_env"
	"github.com/pixality-inc/golang-core/control_flow"
	"github.com/pixality-inc/golang-core/env"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/util"

	"stagen/internal/build"
	"stagen/internal/cli"
	"stagen/internal/config"
	"stagen/pkg/git"
	"stagen/pkg/stagen"
)

type Wiring struct {
	ControlFlow control_flow.ControlFlow
	EnvConfig   *config.Config
	BaseEnv     base_env.BaseEnv
	Log         logger.Logger
	Git         git.Git
	Stagen      stagen.Stagen
	Cli         cli.Cli
}

func New() *Wiring {
	appEnv := env.New(
		"dev",
		util.Ternary(build.CiPipelineId == build.DefaultValue, "", build.CiPipelineId),
		util.Ternary(build.GitTag == build.DefaultValue, "", build.GitTag),
		util.Ternary(build.GitBranch == build.DefaultValue, "", build.GitBranch),
		util.Ternary(build.GitCommit == build.DefaultValue, "", build.GitCommit),
		util.Ternary(build.GitCommitShort == build.DefaultValue, "", build.GitCommitShort),
		time.Now(),
	)

	envConfig, err := config.NewConfigFromEnv()
	if err != nil {
		logger.GetLoggerWithoutContext().WithError(err).Error("Error loading env config")
	}

	baseEnv := base_env.NewBaseEnv(appEnv, logger.DefaultConfig)

	log := baseEnv.Logger()

	controlFlow := control_flow.NewControlFlow(context.Background())

	// Git

	gitTool := git.New()

	// Stagen

	stagenTool := stagen.New(gitTool)

	// Cli

	cliTool := cli.New(stagenTool)

	// Wire

	return &Wiring{
		ControlFlow: controlFlow,
		EnvConfig:   envConfig,
		BaseEnv:     baseEnv,
		Log:         log,
		Git:         gitTool,
		Stagen:      stagenTool,
		Cli:         cliTool,
	}
}

func (w *Wiring) Shutdown() {
	w.ControlFlow.Shutdown()
}
