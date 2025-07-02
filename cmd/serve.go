package cmd

import (
	"fmt"

	"github.com/dskart/particles/api"
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
		api := api.NewAPI(rootLogger)

		port, _ := cmd.Flags().GetInt("port")
		rootLogger.Info().Msgf("serving on port %d", port)
		return ssh.ListenAndServe(fmt.Sprintf(":%d", port), api.Handler)
	},
}
