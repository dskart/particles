package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dskart/particles/api"
	"github.com/dskart/particles/app"
	"github.com/dskart/particles/pkg/shutdown"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 2222, "port to listen on")
	serveCmd.Flags().Int("http-port", 8080, "port for health checks and SSH over websocket (/ssh), 0 to disable")
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
		eg.Go(func() error {
			if err := apiInstance.Server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
				return err
			}
			return nil
		})

		httpPort, _ := cmd.Flags().GetInt("http-port")
		var httpServer *http.Server
		if httpPort != 0 {
			rootLogger.Info().Msgf("serving http on port %d", httpPort)
			httpServer = &http.Server{Addr: fmt.Sprintf(":%d", httpPort), Handler: api.NewHTTPHandler(apiInstance)}
			eg.Go(func() error {
				if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				return nil
			})
		}

		eg.Go(func() error {
			<-ctx.Done()
			_ = apiInstance.Server.Close()
			if httpServer != nil {
				_ = httpServer.Close()
			}
			return nil
		})

		return eg.Wait()
	},
}
