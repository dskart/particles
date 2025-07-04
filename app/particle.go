package app

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Particle struct {
	x, y   float64
	vx, vy float64
	life   int
	color  tcell.Color
}

type ParticleBuffer struct {
	particles []Particle
	mu        sync.RWMutex
}

func NewParticleBuffer() *ParticleBuffer {
	return &ParticleBuffer{
		particles: make([]Particle, 0),
	}
}

func (pb *ParticleBuffer) AddParticle(x, y float64, color tcell.Color) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.particles = append(pb.particles, Particle{
		x:     x,
		y:     y,
		vx:    (rand.Float64() - 0.5) * 20,
		vy:    (rand.Float64() - 0.5) * 20,
		life:  100,
		color: color,
	})
}

func (pb *ParticleBuffer) Length() int {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return len(pb.particles)
}

func (pb *ParticleBuffer) GetParticle(index int) Particle {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return pb.particles[index]
}

func (pb *ParticleBuffer) GetAllParticles() []Particle {
	pb.mu.RLock()
	defer pb.mu.RUnlock()

	dest := make([]Particle, len(pb.particles))
	copy(dest, pb.particles)
	return dest
}

func (pb *ParticleBuffer) UpdateParticle(index int, p Particle) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if index >= len(pb.particles) {
		return fmt.Errorf("%d index out of particle buffer, cannot update", index)
	}
	pb.particles[index] = p
	return nil
}

func (pb *ParticleBuffer) RemoveParticle(index int) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if index >= len(pb.particles) {
		return fmt.Errorf("%d index out of particle buffer, cannot remove", index)
	}
	pb.particles = append(pb.particles[:index], pb.particles[index+1:]...)
	return nil
}
