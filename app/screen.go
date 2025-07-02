package app

import (
	"bytes"
	"context"
	"os"
	"sync"

	"golang.org/x/term"
)

// Keep track of dirty and not dirty cells
type Cell struct {
	ch    rune
	color string
}

type Screen struct {
	width      int
	height     int
	cellBuffer []Cell
	lock       sync.RWMutex
	buf        bytes.Buffer
	tty        *os.File
}

func NewScreen(ctx context.Context) (*Screen, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}

	cellBuffer := make([]Cell, width*height)
	for i := range cellBuffer {
		cellBuffer[i] = Cell{ch: ' ', color: "black"}
	}

	return &Screen{
		width:      width,
		height:     height,
		cellBuffer: cellBuffer,
		tty:        os.Stdout,
	}, nil
}

func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

func (s *Screen) Init() error {
	if err := s.hideCursor(false); err != nil {
		return err
	}
	if err := s.clearScreen(false); err != nil {
		return err
	}
	return nil
}

func (s *Screen) Shutdown() error {
	// if err := s.tty.Sync(); err != nil {
	// 	return err
	// }
	if err := s.showCursor(false); err != nil {
		return err
	}
	if err := s.clearScreen(false); err != nil {
		return err
	}
	return nil
}

func (s *Screen) Clear() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i := range s.cellBuffer {
		s.cellBuffer[i] = Cell{ch: ' ', color: "white"}
	}
	return nil
}

func (s *Screen) SetContent(x, y int, ch rune, color string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cellBuffer[y*s.width+x] = Cell{ch, color}
}

func (s *Screen) Show(ctx context.Context) error {
	if err := s.clearScreen(false); err != nil {
		return err
	}
	s.buf.Reset()

	s.lock.RLock()
	defer s.lock.RUnlock()
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if err := s.drawCell(x, y); err != nil {
				return err
			}
		}
	}

	if _, err := s.buf.WriteTo(s.tty); err != nil {
		return err
	}

	return nil
}

func (s *Screen) drawCell(x, y int) error {
	cell := s.cellBuffer[y*s.width+x]
	if err := s.writeString(string(cell.ch), true); err != nil {
		return err
	}
	return nil
}

func (s *Screen) clearScreen(buf bool) error {
	return s.writeString("\033[2J\033[H", buf)
}

func (s *Screen) hideCursor(buf bool) error {
	return s.writeString("\033[?25l", buf)
}

func (s *Screen) showCursor(buf bool) error {
	return s.writeString("\033[?25h", buf)
}

func (s *Screen) writeString(str string, buf bool) error {
	if buf {
		s.buf.WriteString(str)
	} else {
		_, err := s.tty.WriteString(str)
		return err
	}
	return nil
}
