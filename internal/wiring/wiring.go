package wiring

import (
	"context"

	"github.com/pixality-inc/golang-core/control_flow"
	"github.com/pixality-inc/golang-core/logger"

	"stagen/internal/cli"
	"stagen/internal/config"
	"stagen/pkg/git"
)

type Wiring struct {
	ControlFlow control_flow.ControlFlow
	EnvConfig   *config.Config
	Log         logger.Logger
	Git         git.Git
	Cli         cli.Cli
}

func New() *Wiring {
	// appEnv := env.New(
	// 	 "dev",
	// 	 util.Ternary(build.CiPipelineId == build.DefaultValue, "", build.CiPipelineId),
	// 	 util.Ternary(build.GitTag == build.DefaultValue, "", build.GitTag),
	// 	 util.Ternary(build.GitBranch == build.DefaultValue, "", build.GitBranch),
	// 	 util.Ternary(build.GitCommit == build.DefaultValue, "", build.GitCommit),
	// 	 util.Ternary(build.GitCommitShort == build.DefaultValue, "", build.GitCommitShort),
	// 	 time.Now(),
	// )
	envConfig, err := config.NewConfigFromEnv()
	if err != nil {
		logger.GetLoggerWithoutContext().WithError(err).Error("Error loading env config")
	}

	log := logger.New(logger.DefaultConfig)

	if err := logger.InitLogSpawner(log); err != nil {
		log.WithError(err).Fatal("error initializing log spawner")
	}

	controlFlow := control_flow.NewControlFlow(context.Background())

	// Git

	gitTool := git.New()

	// Cli

	cliTool := cli.New(gitTool)

	// Wire

	return &Wiring{
		ControlFlow: controlFlow,
		EnvConfig:   envConfig,
		Log:         log,
		Git:         gitTool,
		Cli:         cliTool,
	}
}

func (w *Wiring) Shutdown() {
	w.ControlFlow.Shutdown()
}
