package wiring

import (
	"context"
	"time"

	"stagen/internal/build"
	"stagen/internal/config"
	"stagen/pkg/stagen"

	"github.com/pixality-inc/golang-core/base_env"
	"github.com/pixality-inc/golang-core/control_flow"
	"github.com/pixality-inc/golang-core/env"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/util"
)

type Wiring struct {
	ControlFlow control_flow.ControlFlow
	Config      *config.Config
	Log         logger.Logger
	Stagen      stagen.Stagen
}

func New() *Wiring {
	controlFlow := control_flow.NewControlFlow(context.Background())

	cfg := config.LoadConfig()

	appEnv := env.New(
		"dev",
		util.Ternary(build.CiPipelineId == build.DefaultValue, "", build.CiPipelineId),
		util.Ternary(build.GitTag == build.DefaultValue, "", build.GitTag),
		util.Ternary(build.GitBranch == build.DefaultValue, "", build.GitBranch),
		util.Ternary(build.GitCommit == build.DefaultValue, "", build.GitCommit),
		util.Ternary(build.GitCommitShort == build.DefaultValue, "", build.GitCommitShort),
		time.Now(),
	)

	baseEnv := base_env.NewBaseEnv(appEnv, &cfg.Logger)

	log := baseEnv.Logger()

	// Stagen

	stagenEngine := stagen.New(
		&cfg.Stagen,
		&cfg.Site,
	)

	// Wire

	return &Wiring{
		ControlFlow: controlFlow,
		Config:      cfg,
		Log:         log,
		Stagen:      stagenEngine,
	}
}

func (w *Wiring) Shutdown() {
	w.ControlFlow.Shutdown()
}
