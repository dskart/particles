package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
)

func setupScreen(s ssh.Session) (tcell.Screen, error) {
	sshTty, err := NewSSHTty(s)
	if err != nil {
		return nil, err
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	ti, err := sshTty.GetTerminfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get terminfo: %w", err)
	}
	screen, err := tcell.NewTerminfoScreenFromTtyTerminfo(sshTty, ti)
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to init screen: %w", err)
	}

	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	return screen, nil
}

func closeScreen(screen tcell.Screen) func() {
	return func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
}
