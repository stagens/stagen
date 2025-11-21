package main

import (
	"stagen/internal/wiring"
	"stagen/pkg/stagen"
)

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	stagenEngine := stagen.New(
		&wire.Config.Stagen,
		&wire.Config.Site,
	)

	if err := stagenEngine.Build(wire.ControlFlow.Context()); err != nil {
		wire.Log.WithError(err).Fatal("Failed to build stagen website")
	}
}
