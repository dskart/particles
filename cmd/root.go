package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/dskart/particles/pkg/config"
	"github.com/dskart/particles/pkg/logger"
	"github.com/dskart/particles/pkg/shutdown"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

//go:embed banner.txt
var bannerArt string

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "make output more verbose")
	rootCmd.PersistentFlags().StringP("config", "c", "", "read configuration from this file")
}

const (
	cfgEnvPrefix = "PARTICLES"
	serviceName  = "particles"
)

var rootConfig Config
var rootLogger *zerolog.Logger

var rootCmd = &cobra.Command{
	Use:           filepath.Base(os.Args[0]),
	SilenceErrors: false,
	SilenceUsage:  false,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		rootLogger = logger.NewLogger(verbose, serviceName)

		configFilePath, _ := cmd.Flags().GetString("config")
		cfgOpts := []func(*config.UnmarshalConfigOptions){config.WithPrefix(cfgEnvPrefix)}
		if configFilePath != "" {
			cfgOpts = append(cfgOpts, config.WithFilePath(configFilePath))
		}

		if err := config.UnmarshalConfig(context.Background(), &rootConfig, cfgOpts...); err != nil {
			return err
		}

		// wait forever for sig signal
		go func() {
			WaitForTermSignal()
		}()

		fmt.Println(bannerArt)

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// Do a graceful shutdown
		shutdown.Shutdown()
		return nil
	},
}

func WaitForTermSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	sig := <-sigs
	rootLogger.Error().Str("signal", sig.String()).Msg("received signal, shutting down")

	// Do a graceful shutdown
	shutdown.Shutdown()
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootLogger.Fatal().Stack().Err(err).Msg(err.Error())
	}
}
