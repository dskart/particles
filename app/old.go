package app

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"golang.org/x/term"
// )

// func (a *App) RunOld(ctx context.Context) error {
// 	// Create log channel for particle logger
// 	logChan := make(chan string, 1000)

// 	particleLogger := NewParticleLogger(ctx, a.Logger.GetLevel(), logChan)

// 	logCounter := 0
// 	go func() {
// 		for {
// 			time.Sleep(1 * time.Second)
// 			particleLogger.Info().Msg(fmt.Sprintf("Loop: %d", logCounter))
// 			logCounter++
// 		}
// 	}()

// 	// Clear terminal and set black background
// 	fmt.Print("\033[2J")  // Clear screen
// 	fmt.Print("\033[40m") // Set background to black
// 	fmt.Print("\033[37m") // Set foreground to white for visibility

// 	counter := 0
// 	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
// 	defer ticker.Stop()

// 	// Buffer for recent logs (keep last 5 logs)
// 	recentLogs := make([]string, 0, a.Config.MaxNLogs)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Print("\033[0m") // Reset colors
// 			return nil
// 		case logMsg := <-logChan:
// 			logMsg = strings.TrimSuffix(logMsg, "\n")
// 			logMsg = strings.TrimSpace(logMsg)
// 			recentLogs = append(recentLogs, logMsg)
// 			if len(recentLogs) > a.Config.MaxNLogs {
// 				recentLogs = recentLogs[1:]
// 			}

// 		case <-ticker.C:
// 			// Get current terminal width dynamically
// 			width, _, err := term.GetSize(int(os.Stdout.Fd()))
// 			if err != nil {
// 				width = 80 // fallback width
// 			}

// 			// Clear screen and move cursor to top left
// 			fmt.Print("\033[2J\033[H")

// 			// Print logs in top left
// 			if len(recentLogs) > 0 {
// 				for _, log := range recentLogs {
// 					fmt.Printf("%s\n", log)
// 				}
// 			}

// 			// Print counter in top right
// 			fmt.Printf("\033[H\033[%dCLoop: %d", width-15, counter)
// 			counter++
// 		}
// 	}
// }
