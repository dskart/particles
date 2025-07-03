package app

import "github.com/gdamore/tcell/v2"

type VirtualScreen struct {
	realScreen tcell.Screen
	offsetX    int
	offsetY    int
	maxWidth   int
	maxHeight  int
	viewWidth  int
	viewHeight int
}

// SetContent implements app.Renderer.SetContent with offset and bounds checking
func (vs *VirtualScreen) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) {
	// Check bounds within the virtual screen
	if x >= 0 && x < vs.maxWidth && y >= 0 && y < vs.maxHeight {
		realX := x + vs.offsetX
		realY := y + vs.offsetY
		// Check bounds within the real screen view
		if realX >= 0 && realX < vs.viewWidth && realY >= 0 && realY < vs.viewHeight {
			vs.realScreen.SetContent(realX, realY, mainc, combc, style)
		}
	}
}

// Size implements app.Renderer.Size - returns the virtual screen size (simulation area size)
func (vs *VirtualScreen) Size() (int, int) {
	return vs.maxWidth, vs.maxHeight
}

// Clear implements app.Renderer.Clear - clears the virtual screen area
func (vs *VirtualScreen) Clear() {
	for y := 0; y < vs.maxHeight; y++ {
		for x := 0; x < vs.maxWidth; x++ {
			vs.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}
}
