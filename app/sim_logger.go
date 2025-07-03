package app

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
)

type simLoggerWriter struct {
	logs     *[]string
	maxNLogs int
}

func (slw *simLoggerWriter) Write(p []byte) (n int, err error) {
	*slw.logs = append(*slw.logs, string(p))

	if len(*slw.logs) > slw.maxNLogs {
		*slw.logs = (*slw.logs)[1:]
	}
	return len(p), nil
}

type SimLogger struct {
	logs   *[]string
	logger *zerolog.Logger
}

func NewSimLogger(ctx context.Context, level zerolog.Level, maxNLogs int) *SimLogger {
	logSlice := make([]string, 0, maxNLogs)
	channelWriter := &simLoggerWriter{logs: &logSlice, maxNLogs: maxNLogs}

	simLogger := zerolog.New(channelWriter).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &SimLogger{
		logs:   &logSlice,
		logger: &simLogger,
	}
}

func (l *SimLogger) Logger() *zerolog.Logger {
	return l.logger
}

func (l *SimLogger) Render(renderer Renderer) {
	for y, log := range *l.logs {
		for x, r := range log {
			renderer.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
}
