package main

import (
	"context"

	"github.com/pixality-inc/golang-core/logger"
	"github.com/spf13/cobra"

	"stagen/internal/cli"
	"stagen/internal/config"
	"stagen/internal/wiring"
)

func runCommand(rootCtx context.Context, cli cli.Cli) {
	log := logger.GetLogger(rootCtx)

	// Root
	rootCmd := &cobra.Command{
		Use:   "stagen",
		Short: "Static website generator",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.WithError(err).Fatal()
			}
		},
	}

	// Init

	{
		initCmd := &cobra.Command{
			Use:   "init dir",
			Short: "Initialize new project in directory",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
				workDir := args[0]
				name := cmd.Flag("name").Value.String()

				if err := cli.Init(cmd.Context(), workDir, name); err != nil {
					log.WithError(err).Fatal()
				}
			},
		}

		initCmd.Flags().StringP("name", "n", "Stagen Website", "website name")

		rootCmd.AddCommand(initCmd)
	}

	// Build

	{
		buildCmd := &cobra.Command{
			Use:   "build [dir]",
			Short: "Build project in directory [dir]",
			Args:  cobra.MaximumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
				workDir := config.RootDir()

				if len(args) > 0 {
					workDir = args[0]
				}

				if err := cli.Build(cmd.Context(), workDir); err != nil {
					log.WithError(err).Fatal()
				}
			},
		}

		rootCmd.AddCommand(buildCmd)
	}

	// Execute

	if err := rootCmd.ExecuteContext(rootCtx); err != nil {
		log.WithError(err).Fatal("failed to run command")
	}
}

func main() {
	wire := wiring.New()
	defer wire.Shutdown()

	runCommand(wire.ControlFlow.Context(), wire.Cli)
}
