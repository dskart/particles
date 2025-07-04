package cmd

import (
	"context"
	"fmt"

	"github.com/dskart/particles/api"
	"github.com/dskart/particles/app"
	"github.com/dskart/particles/pkg/shutdown"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 2222, "port to listen on")
}

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(cancel)
		rootLogger.Info().Msg("serve command")
		appInstance, err := app.NewApp(rootLogger, rootConfig.App)
		if err != nil {
			return fmt.Errorf("failed to create APP: %w", err)
		}

		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			return appInstance.Run(ctx)
		})

		apiInstance, err := api.NewAPI(rootLogger, rootConfig.API, appInstance)
		if err != nil {
			return fmt.Errorf("failed to create API: %w", err)
		}

		port, _ := cmd.Flags().GetInt("port")
		rootLogger.Info().Msgf("serving on port %d", port)
		apiInstance.Server.Addr = fmt.Sprintf(":%d", port)
		return apiInstance.Server.ListenAndServe()
	},
}
