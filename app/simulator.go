package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
)

type SimConfig struct {
	Width   int `yaml:"Width" env:"WIDTH"`
	Height  int `yaml:"Height" env:"HEIGHT"`
	Gravity int `yaml:"Gravity" env:"GRAVITY"`
}

func NewSimConfig() SimConfig {
	return SimConfig{
		Width:   80,
		Height:  25,
		Gravity: 1,
	}
}

func (c *SimConfig) Validate() error {
	if c.Width < 80 {
		return fmt.Errorf("width too small %d", c.Width)
	}
	if c.Height < 25 {
		return fmt.Errorf("height too small %d", c.Width)
	}
	return nil
}

type Simulator struct {
	config         SimConfig
	particleBuffer *ParticleBuffer
	width          int
	height         int
	startTime      time.Time
}

func NewSimulator(ctx context.Context, level zerolog.Level, simConfig SimConfig) (*Simulator, error) {
	return &Simulator{
		particleBuffer: NewParticleBuffer(),
		config:         simConfig,
		width:          simConfig.Width * 2,  // Each character can hold 2 horizontal pixels
		height:         simConfig.Height * 2, // Each character can hold 2 vertical pixels
		startTime:      time.Now(),
	}, nil
}

func (s *Simulator) ElapsedTime() time.Duration {
	return time.Since(s.startTime)
}

func (s *Simulator) Size() (int, int) {
	return s.config.Width, s.config.Height
}

func (s *Simulator) AddParticle(x, y float64, color tcell.Color) {
	s.particleBuffer.AddParticle(x, y, color)
}

func (s *Simulator) NumOfParticles() int {
	return s.particleBuffer.Length()
}

func (s *Simulator) Run(ctx context.Context, logger *zerolog.Logger) error {
	delta := 50 * time.Millisecond
	ticker := time.NewTicker(delta)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			logger.Debug().TimeDiff("simTime", time.Now(), s.startTime).Msg("Sim Loop")
			if err := s.Update(delta); err != nil {
				return fmt.Errorf("could not Update: %w", err)
			}
		}
	}

}

func (s *Simulator) Update(delta time.Duration) error {
	numParticles := s.particleBuffer.Length()
	// Going backwards so that even if new particles are added we don't break
	for i := numParticles - 1; i >= 0; i-- {
		p := s.particleBuffer.GetParticle(i)

		p.x += p.vx * delta.Seconds()
		p.y += p.vy * delta.Seconds()
		p.vy += float64(s.config.Gravity)
		p.life--

		// Remove dead or out-of-bounds particles
		if p.life <= 0 || p.x < 0 || p.x >= float64(s.width) ||
			p.y < 0 || p.y >= float64(s.height) {
			if err := s.particleBuffer.RemoveParticle(i); err != nil {
				return fmt.Errorf("could not remove particle %d: %w", i, err)
			}
		} else {
			if err := s.particleBuffer.UpdateParticle(i, p); err != nil {
				return fmt.Errorf("could not update particle %d: %w", i, err)
			}
		}
	}
	return nil
}

func (s *Simulator) Render(renderer Renderer) {
	// Create a 2D array to track which sub-pixels are set and their colors
	pixels := make([][]bool, s.height)
	colors := make([][]tcell.Color, s.height)
	for i := range pixels {
		pixels[i] = make([]bool, s.width)
		colors[i] = make([]tcell.Color, s.width)
	}

	// Set pixels for particles
	particles := s.particleBuffer.GetAllParticles()
	for _, p := range particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < s.width && y >= 0 && y < s.height {
			pixels[y][x] = true
			colors[y][x] = p.color
		}
	}

	// Convert pixel array to Unicode block characters
	for y := 0; y < s.height; y += 2 {
		for x := 0; x < s.width; x += 2 {
			char, color := getBlockCharWithColor(
				pixels, colors, x, y, s.width, s.height,
			)
			if char != ' ' {
				style := tcell.StyleDefault.Foreground(color)
				renderer.SetContent(x/2, y/2, char, nil, style)
			}
		}
	}
}

func getBlockCharWithColor(pixels [][]bool, colors [][]tcell.Color, x, y, width, height int) (rune, tcell.Color) {
	// Unicode block characters for 2x2 pixel representation
	var mask int
	color := tcell.ColorWhite // default color

	if y < height && x < width && pixels[y][x] {
		mask |= 1
		color = colors[y][x]
	}
	if y < height && x+1 < width && pixels[y][x+1] {
		mask |= 2
		if mask == 2 { // first pixel found
			color = colors[y][x+1]
		}
	}
	if y+1 < height && x < width && pixels[y+1][x] {
		mask |= 4
		if mask == 4 { // first pixel found
			color = colors[y+1][x]
		}
	}
	if y+1 < height && x+1 < width && pixels[y+1][x+1] {
		mask |= 8
		if mask == 8 { // first pixel found
			color = colors[y+1][x+1]
		}
	}

	blockChars := []rune{
		' ', // 0000
		'▘', // 0001
		'▝', // 0010
		'▀', // 0011
		'▖', // 0100
		'▌', // 0101
		'▞', // 0110
		'▛', // 0111
		'▗', // 1000
		'▚', // 1001
		'▐', // 1010
		'▜', // 1011
		'▄', // 1100
		'▙', // 1101
		'▟', // 1110
		'█', // 1111
	}

	return blockChars[mask], color
}
