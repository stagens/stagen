package main

import (
	"stagen/internal/wiring"
)

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	// @todo build command
	// @todo dev command (web server with building of-the-fly)
	// @todo init command (create directories structure, clone template, create config, Makefile, .gitignore)

	if err := wire.Stagen.Build(wire.ControlFlow.Context()); err != nil {
		wire.Log.WithError(err).Fatal("Failed to build stagen website")
	}
}
