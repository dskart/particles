package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func NewLogger(verbose bool, serviceName string) *zerolog.Logger {
	lvl := zerolog.InfoLevel
	if verbose {
		lvl = zerolog.DebugLevel
	}

	var w io.Writer
	w = os.Stdout
	if term.IsTerminal(unix.Stdout) {
		w = &zerolog.ConsoleWriter{Out: w}
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log := zerolog.New(w).Level(lvl).With().Timestamp().Caller().Str("service", serviceName).Logger()

	return &log
}
