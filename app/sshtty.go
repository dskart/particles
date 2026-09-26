package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/terminfo"
	"github.com/gliderlabs/ssh"
)

// SSHTty implements tcell.Tty interface for SSH sessions
type SSHTty struct {
	session  ssh.Session
	winCh    <-chan ssh.Window
	ptyReq   ssh.Pty
	resizeCb func()
	stopCh   chan struct{}
	started  bool
}

func NewSSHTty(s ssh.Session) (*SSHTty, error) {
	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		return nil, fmt.Errorf("no PTY requested")
	}

	return &SSHTty{
		session: s,
		winCh:   winCh,
		ptyReq:  ptyReq,
		stopCh:  make(chan struct{}),
	}, nil
}

func (s *SSHTty) GetTerminfo() (*terminfo.Terminfo, error) {
	if s.ptyReq.Term == "" {
		return nil, fmt.Errorf("no ptyReq.Term")
	}
	ti, err := tcell.LookupTerminfo(s.ptyReq.Term)
	if err != nil {
		return nil, nil
	}
	return ti, nil
}

func (s *SSHTty) Start() error {
	if s.started {
		return nil
	}
	s.started = true

	go func() {
		for {
			select {
			case win := <-s.winCh:
				s.ptyReq.Window = win
				if s.resizeCb != nil {
					s.resizeCb()
				}
			case <-s.stopCh:
				return
			}
		}
	}()

	return nil
}

func (s *SSHTty) Stop() error {
	if !s.started {
		return nil
	}
	s.started = false
	close(s.stopCh)
	s.stopCh = make(chan struct{}) // Reset for potential future Start()
	return nil
}

func (s *SSHTty) Drain() error {
	// SSH session doesn't need draining like /dev/tty
	return nil
}

func (s *SSHTty) NotifyResize(cb func()) {
	s.resizeCb = cb
}

func (s *SSHTty) WindowSize() (tcell.WindowSize, error) {
	return tcell.WindowSize{
		Width:       s.ptyReq.Window.Width,
		Height:      s.ptyReq.Window.Height,
		PixelWidth:  0, // SSH doesn't provide pixel dimensions
		PixelHeight: 0,
	}, nil
}

func (s *SSHTty) Write(b []byte) (int, error) {
	return s.session.Write(b)
}

func (s *SSHTty) Read(b []byte) (int, error) {
	return s.session.Read(b)
}

func (s *SSHTty) Close() error {
	if s.started {
		_ = s.Stop()
	}
	return s.session.Close()
}
