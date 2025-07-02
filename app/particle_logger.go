package app

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
)

type particleLoggerWriter struct {
	logs     *[]string
	maxNLogs int
}

func (plw *particleLoggerWriter) Write(p []byte) (n int, err error) {
	*plw.logs = append(*plw.logs, string(p))

	if len(*plw.logs) > plw.maxNLogs {
		*plw.logs = (*plw.logs)[1:]
	}
	return len(p), nil
}

type ParticleLogger struct {
	logs   *[]string
	logger *zerolog.Logger
}

func NewParticleLogger(ctx context.Context, level zerolog.Level, maxNLogs int) *ParticleLogger {
	logSlice := make([]string, 0, maxNLogs)
	channelWriter := &particleLoggerWriter{logs: &logSlice, maxNLogs: maxNLogs}

	particleLogger := zerolog.New(channelWriter).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &ParticleLogger{
		logs:   &logSlice,
		logger: &particleLogger,
	}
}

func (l *ParticleLogger) Logger() *zerolog.Logger {
	return l.logger
}

func (l *ParticleLogger) Render(screen tcell.Screen) {
	for y, log := range *l.logs {
		for x, r := range log {
			screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
}
