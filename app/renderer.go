package app

import "github.com/gdamore/tcell/v2"

// Renderer is a minimal interface for rendering content to a screen
type Renderer interface {
	SetContent(x, y int, mainc rune, combc []rune, style tcell.Style)
}
