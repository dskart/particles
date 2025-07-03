package api

import (
	"io"
	"sync/atomic"

	"github.com/dskart/particles/app"
	"github.com/gliderlabs/ssh"
	"github.com/rs/zerolog"
)

type API struct {
	Logger  *zerolog.Logger
	Handler ssh.Handler
	Config  Config
}

func NewAPI(logger *zerolog.Logger, config Config, app *app.App) (*API, error) {
	config.Validate()
	api := &API{Logger: logger, Config: config}

	var numActiveSessions atomic.Int32

	h := ssh.Handler(func(s ssh.Session) {
		sessLogger := logger.With().Str("remoteAddr", s.RemoteAddr().String()).Str("user", s.User()).Logger()
		if numSessions := numActiveSessions.Load(); numSessions >= int32(config.MaxNumSessions) {
			sessLogger.Warn().Int("numSessions", config.MaxNumSessions).Msgf("Too many sessions %d/%d", numSessions, config.MaxNumSessions)
			io.WriteString(s, "Too many sessions")
			return
		}

		currentNumSessions := numActiveSessions.Add(1)
		defer func() {
			currentNumSessions = numActiveSessions.Add(-1)
			sessLogger.Info().Int32("numSessions", numActiveSessions.Load()).Msg("SSH session ended")

		}()

		sessLogger.Info().Int32("numSessions", currentNumSessions).Msg("SSH session started")
		if err := app.HandleSSHSession(s, &sessLogger, &numActiveSessions, config.MaxNumSessions); err != nil {
			sessLogger.Err(err).Msg("")
			io.WriteString(s, err.Error())
			s.Exit(1)
			return
		}

		s.Exit(0)
	})

	api.Handler = h
	return api, nil
}
