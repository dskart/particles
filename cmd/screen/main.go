package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dskart/particles/app"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Handle Ctrl+C
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		cancel()
	}()

	screen, err := app.NewScreen(ctx)
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}

	// hello := "Hello, World!"
	// for x, r := range hello {
	// 	screen.SetContent(x, 0, r, "white")
	// }

	width, height := screen.Size()
	sim, err := app.NewSimulator(width, height)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		sim.AddParticle(float64(rand.Intn(width)), float64(rand.Intn(height)))
	}

	// counter := 0
	ticker := time.NewTicker(50 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			screen.Shutdown()
			return
		case <-ticker.C:
			screen.Clear()
			// fmt.Println("foo")
			// hello := fmt.Sprintf("counter: %s ", strconv.Itoa(counter))
			// for x, r := range hello {
			// 	screen.SetContent(x, 0, r, "white")
			// 	counter++
			// }

			sim.Render2(screen)
			screen.Show(ctx)
		}
	}
}
