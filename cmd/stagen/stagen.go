package main

import (
	"context"

	"github.com/pixality-inc/golang-core/logger"
	"github.com/spf13/cobra"

	"stagen/internal/cli"
	"stagen/internal/config"
	"stagen/internal/wiring"
)

func runCommand(rootCtx context.Context, cliTool cli.Cli) {
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
		cmd := &cobra.Command{
			Use:   "init dir",
			Short: "Initialize new project in directory",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
				workDir := args[0]
				name := cmd.Flag("name").Value.String()

				if err := cliTool.Init(cmd.Context(), workDir, name, true); err != nil {
					log.WithError(err).Fatal()
				}
			},
		}

		cmd.Flags().StringP("name", "n", "Stagen Website", "website name")

		rootCmd.AddCommand(cmd)
	}

	// Build

	rootCmd.AddCommand(&cobra.Command{
		Use:   "build [dir]",
		Short: "Build project in directory [dir]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
			workDir := config.RootDir()

			if len(args) > 0 {
				workDir = args[0]
			}

			if err := cliTool.Build(cmd.Context(), workDir); err != nil {
				log.WithError(err).Fatal()
			}
		},
	})

	// Watch

	rootCmd.AddCommand(&cobra.Command{
		Use:   "watch [dir]",
		Short: "Watch project and rebuild in directory [dir]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
			workDir := config.RootDir()

			if len(args) > 0 {
				workDir = args[0]
			}

			if err := cliTool.Watch(cmd.Context(), workDir); err != nil {
				log.WithError(err).Fatal()
			}
		},
	})

	// Web

	rootCmd.AddCommand(&cobra.Command{
		Use:   "web [dir]",
		Short: "Serve project over http in directory [dir]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
			workDir := config.RootDir()

			if len(args) > 0 {
				workDir = args[0]
			}

			if err := cliTool.Web(cmd.Context(), workDir); err != nil {
				log.WithError(err).Fatal()
			}
		},
	})

	// Dev

	rootCmd.AddCommand(&cobra.Command{
		Use:   "dev [dir]",
		Short: "Serve project over http and watch for changes in directory [dir]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) { //nolint:contextcheck
			workDir := config.RootDir()

			if len(args) > 0 {
				workDir = args[0]
			}

			if err := cliTool.Dev(cmd.Context(), workDir); err != nil {
				log.WithError(err).Fatal()
			}
		},
	})

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
