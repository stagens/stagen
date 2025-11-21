package main

import "stagen/internal/wiring"

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	wire.Log.Info("Hello, world!")
}
