package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(noopCmd)
}

var noopCmd = &cobra.Command{
	Use: "noop",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootLogger.Info().Msg("noop command")
		return fmt.Errorf("test error")
	},
}
