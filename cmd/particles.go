package cmd

import (
	"context"

	"github.com/dskart/particles/app"
	"github.com/dskart/particles/pkg/shutdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(particlesCmd)
}

var particlesCmd = &cobra.Command{
	Use: "particles",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(func() {
			cancel()
		})

		app := app.NewApp(rootLogger, rootConfig.App)
		if err := app.Run(ctx); err != nil {
			return err
		}

		return nil
	},
}
