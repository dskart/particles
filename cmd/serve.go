package cmd

import (
	"fmt"

	"github.com/dskart/particles/api"
	"github.com/dskart/particles/app"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 2222, "port to listen on")
}

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootLogger.Info().Msg("serve command")
		appInstance, err := app.NewApp(rootLogger, rootConfig.App)
		if err != nil {
			return fmt.Errorf("failed to create APP: %w", err)
		}

		apiInstance, err := api.NewAPI(rootLogger, rootConfig.API, appInstance)
		if err != nil {
			return fmt.Errorf("failed to create API: %w", err)
		}

		port, _ := cmd.Flags().GetInt("port")
		rootLogger.Info().Msgf("serving on port %d", port)
		return ssh.ListenAndServe(fmt.Sprintf(":%d", port), apiInstance.Handler)
	},
}
