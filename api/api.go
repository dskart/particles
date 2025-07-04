package api

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/dskart/particles/app"
	"github.com/gliderlabs/ssh"
	"github.com/rs/zerolog"
	golangSsh "golang.org/x/crypto/ssh"
)

type API struct {
	Logger  *zerolog.Logger
	Handler ssh.Handler
	Config  Config
	Server  *ssh.Server
}

func NewAPI(logger *zerolog.Logger, config Config, app *app.App) (*API, error) {
	config.Validate()
	api := &API{Logger: logger, Config: config}
	server := &ssh.Server{}
	if config.SSHHostKey != "" {
		logger.Info().Msg("Using SSH host key from configuration")
		signer, err := parseSSHHostKey(config.SSHHostKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SSH host key: %w", err)
		}
		server.HostSigners = []ssh.Signer{signer}
	} else {
		logger.Warn().Msg("No SSH host key provided in configuration, server will use a generated key")
	}

	api.Server = server

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
	api.Server.Handler = h
	return api, nil
}

// parseSSHHostKey parses an SSH host key from various formats
func parseSSHHostKey(hostKey string) (ssh.Signer, error) {
	// Try to parse as PEM-encoded private key
	block, _ := pem.Decode([]byte(hostKey))
	if block != nil {
		// Parse different key types
		switch block.Type {
		case "RSA PRIVATE KEY":
			privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
			}
			return golangSsh.NewSignerFromKey(privateKey)
		case "PRIVATE KEY":
			privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
			}
			return golangSsh.NewSignerFromKey(privateKey)
		default:
			return nil, fmt.Errorf("unsupported key type: %s", block.Type)
		}
	}

	// If not PEM, try to parse as raw key bytes
	privateKey, err := golangSsh.ParseRawPrivateKey([]byte(hostKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw private key: %w", err)
	}

	return golangSsh.NewSignerFromKey(privateKey)
}
