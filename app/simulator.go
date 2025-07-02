package app

import (
	"math/rand/v2"

	"github.com/gdamore/tcell/v2"
)

type Particle struct {
	x, y   float64
	vx, vy float64
	life   int
}

type Simulator struct {
	particles []Particle
	width     int
	height    int
}

func NewSimulator(width, height int) (*Simulator, error) {
	return &Simulator{
		particles: make([]Particle, 0),
		width:     width * 2,  // Each character can hold 2 horizontal pixels
		height:    height * 2, // Each character can hold 2 vertical pixels
	}, nil
}

func (s *Simulator) AddParticle(x, y float64) {
	s.particles = append(s.particles, Particle{
		x:    x,
		y:    y,
		vx:   (rand.Float64() - 0.5) * 20,
		vy:   (rand.Float64() - 0.5) * 20,
		life: 100,
	})
}

func (s *Simulator) Update() {
	// Update particles
	for i := len(s.particles) - 1; i >= 0; i-- {
		p := &s.particles[i]

		p.x += p.vx * 0.1
		p.y += p.vy * 0.1
		p.vy += 0.5 // gravity
		p.life--

		// Remove dead or out-of-bounds particles
		if p.life <= 0 || p.x < 0 || p.x >= float64(s.width) ||
			p.y < 0 || p.y >= float64(s.height) {
			s.particles = append(s.particles[:i], s.particles[i+1:]...)
		}
	}
}

func (s *Simulator) Render(screen tcell.Screen) {
	s.Update()

	// Create a 2D array to track which sub-pixels are set
	pixels := make([][]bool, s.height)
	for i := range pixels {
		pixels[i] = make([]bool, s.width)
	}

	// Set pixels for particles
	for _, p := range s.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < s.width && y >= 0 && y < s.height {
			pixels[y][x] = true
		}
	}

	// Convert pixel array to Unicode block characters
	for y := 0; y < s.height; y += 2 {
		for x := 0; x < s.width; x += 2 {
			char := getBlockChar(
				pixels, x, y, s.width, s.height,
			)
			if char != ' ' {
				screen.SetContent(x/2, y/2, char, nil, tcell.StyleDefault)
			}
		}
	}
}

func (s *Simulator) Render2(screen *Screen) {
	s.Update()

	// Create a 2D array to track which sub-pixels are set
	pixels := make([][]bool, s.height)
	for i := range pixels {
		pixels[i] = make([]bool, s.width)
	}

	// Set pixels for particles
	for _, p := range s.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < s.width && y >= 0 && y < s.height {
			pixels[y][x] = true
		}
	}

	// Convert pixel array to Unicode block characters
	for y := 0; y < s.height; y += 2 {
		for x := 0; x < s.width; x += 2 {
			char := getBlockChar(
				pixels, x, y, s.width, s.height,
			)
			if char != ' ' {
				screen.SetContent(x/2, y/2, char, "")
			}
		}
	}
}

func getBlockChar(pixels [][]bool, x, y, width, height int) rune {
	// Unicode block characters for 2x2 pixel representation
	var mask int

	if y < height && x < width && pixels[y][x] {
		mask |= 1
	}
	if y < height && x+1 < width && pixels[y][x+1] {
		mask |= 2
	}
	if y+1 < height && x < width && pixels[y+1][x] {
		mask |= 4
	}
	if y+1 < height && x+1 < width && pixels[y+1][x+1] {
		mask |= 8
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

	return blockChars[mask]
}
