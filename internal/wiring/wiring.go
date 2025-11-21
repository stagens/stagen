package wiring

import (
	"context"
	"time"

	"stagen/internal/config"

	"github.com/pixality-inc/golang-core/base_env"
	"github.com/pixality-inc/golang-core/control_flow"
	"github.com/pixality-inc/golang-core/env"
	"github.com/pixality-inc/golang-core/logger"
)

type Wiring struct {
	ControlFlow control_flow.ControlFlow
	Config      *config.Config
	Log         logger.Logger
}

func New() *Wiring {
	controlFlow := control_flow.NewControlFlow(context.Background())

	cfg := config.LoadConfig()

	appEnv := env.New("dev", "", "", "", "", "", time.Now())

	baseEnv := base_env.NewBaseEnv(appEnv, &cfg.Logger)

	log := baseEnv.Logger()

	return &Wiring{
		ControlFlow: controlFlow,
		Config:      cfg,
		Log:         log,
	}
}

func (w *Wiring) Shutdown() {
	w.ControlFlow.Shutdown()
}
