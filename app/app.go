package app

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"

	"github.com/rs/zerolog"
)

type App struct {
	config Config
	sim    *Simulator
	logger *zerolog.Logger
}

func NewApp(logger *zerolog.Logger, config Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	sim, err := NewSimulator(context.Background(), logger.GetLevel(), config.SimConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared simulator: %w", err)
	}

	return &App{
		config: config,
		sim:    sim,
		logger: logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info().Msg("Starting sim")
	return a.sim.Run(ctx, a.logger)
}

func (a *App) HandleSSHSession(s ssh.Session, sessLogger *zerolog.Logger, numActiveSessions *atomic.Int32, maxNumSession int) error {
	ctx := s.Context()

	screen, err := setupScreen(s)
	if err != nil {
		return fmt.Errorf("could not setup screen: %w", err)
	}
	defer closeScreen(screen)

	clientColor, hasColor := GetAvailableColor()
	if !hasColor {
		clientColor = getRandColor()
	}
	defer func() {
		if hasColor {
			ReleaseColor(clientColor)
		}
	}()

	maxSessTime, err := time.ParseDuration(a.config.MaxSessTime)
	if err != nil {
		return err
	}

	ok, err := a.showWelcomeMessage(ctx, screen, clientColor, maxSessTime)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	sim := a.sim
	termWidth, termHeight := screen.Size()
	simWidth, simHeight := sim.Size()
	box := newBox(screen, simWidth, simHeight, termWidth, termHeight, clientColor)

	sessionStartTime := time.Now()
	sessEndTime := sessionStartTime.Add(maxSessTime)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if time.Now().After(sessEndTime) {
				return nil
			}
			screen.Clear()

			box.RenderBox(screen, termWidth, termHeight)
			sim.Render(box)

			// Status info
			sessElapsedTime := time.Since(sessionStartTime)
			statusMsg := fmt.Sprintf("Sim Time: %s | Sess Time: %s | Time Left: %s | Color: %s | Users %d/%d | Num Particles %d",
				sim.ElapsedTime().Truncate(time.Second),
				sessElapsedTime.Truncate(time.Second),
				time.Until(sessEndTime).Truncate(time.Second),
				colorName(clientColor),
				numActiveSessions.Load(),
				maxNumSession,
				sim.NumOfParticles(),
			)
			if len(statusMsg) <= termWidth {
				for x, r := range statusMsg {
					screen.SetContent(x, 0, r, nil, tcell.StyleDefault.Foreground(clientColor))
				}
			}

			screen.Show()
		default:
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventResize:
					termWidth, termHeight = screen.Size()
					screen.Clear() // Clear entire screen on resize
					screen.Sync()
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
						return nil
					}
				case *tcell.EventMouse:
					mx, my := ev.Position()
					if ev.Buttons()&(tcell.Button1|tcell.Button2) != 0 {
						if box.InBox(mx, my) {
							simX, simY := box.ToSimCoordinates(mx, my)
							sim.AddParticle(simX, simY, clientColor)
						}
					}
				}
			}
		}
	}

}

// showWelcomeMessage displays instructions to the user and waits for them to press enter
func (a *App) showWelcomeMessage(ctx context.Context, screen tcell.Screen, clientColor tcell.Color, maxSessTime time.Duration) (bool, error) {
	screen.Clear()

	welcomeLines := []string{
		"Welcome to Particles!",
		"",
		"This is a shared particle physics simulation that you can interact with.",
		"",
		fmt.Sprintf("Your particle color: %s", colorName(clientColor)),
		"",
		"How to play:",
		"• Click anywhere in the simulation area to add particles",
		"• Each user gets a unique color for their particles",
		"• Watch as particles interact with gravity and physics",
		fmt.Sprintf("• You will be logged out automatically after %s", maxSessTime.Truncate(time.Second)),
		"• Press Ctrl+C or Escape to exit",
		"",
		"Press Enter to start the simulation...",
	}

	width, height := screen.Size()
	startY := max(0, (height-len(welcomeLines))/2)

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)

	for i, line := range welcomeLines {
		y := startY + i
		if y >= height {
			break
		}

		x := max(0, (width-len(line))/2)

		lineStyle := style
		switch i {
		case 0: // Title line
			lineStyle = titleStyle
		case 4: // Color line
			lineStyle = tcell.StyleDefault.Foreground(clientColor).Bold(true)
		}

		for j, r := range line {
			if x+j < width {
				screen.SetContent(x+j, y, r, nil, lineStyle)
			}
		}
	}

	screen.Show()
	// Wait for user to press Enter
	for {
		select {
		case <-ctx.Done():
			return false, nil
		default:
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEnter {
						return true, nil
					} else if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEscape {
						return false, nil
					}
				}
			}
		}
	}
}
