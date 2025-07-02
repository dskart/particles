package app

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
)

type App struct {
	Logger *zerolog.Logger
	Config Config
}

func NewApp(logger *zerolog.Logger, config Config) *App {
	if err := config.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("Invalid config")
	}

	return &App{
		Logger: logger,
		Config: config,
	}
}

// TODO: create new screen for all ssh connections, keep track of a sim in the back instead
func (a *App) Run(ctx context.Context) error {
	particleLogger := NewParticleLogger(ctx, a.Logger.GetLevel(), a.Config.MaxNLogs)

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen, err := tcell.NewTerminfoScreenFromTty(nil)
	if err != nil {
		return err
	}

	if err := screen.Init(); err != nil {
		return err
	}
	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()
	width, height := screen.Size()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	sim, err := NewSimulator(width, height)
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		sim.AddParticle(float64(rand.Intn(width)), float64(rand.Intn(height)))
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	counter := 0
	for {
		select {
		case <-ticker.C:
			particleLogger.Logger().Debug().Msg(fmt.Sprintf("test %d", counter))
			screen.Clear()
			sim.Render(screen)
			particleLogger.Render(screen)

			counterMsg := fmt.Sprintf("Loop: %d", counter)
			for x, r := range counterMsg {
				screen.SetContent(width-(len(counterMsg)+5)+x, 0, r, nil, tcell.StyleDefault)
			}
			screen.Show()
			counter++
		default:
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventResize:
					width, _ = screen.Size()
					screen.Sync()
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
						return nil
					}
				case *tcell.EventMouse:
					mx, my := ev.Position()
					if ev.Buttons()&(tcell.Button1|tcell.Button2) != 0 {
						sim.AddParticle(float64(mx*2), float64(my*2))
					}

				}
			}
		}
	}
}
