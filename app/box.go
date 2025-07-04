package app

import "github.com/gdamore/tcell/v2"

type Box struct {
	simWidth    int
	simHeight   int
	termWidth   int
	termHeight  int
	boxX        int
	boxY        int
	clientColor tcell.Color
	screen      tcell.Screen
}

func newBox(screen tcell.Screen, simWidth int, simHeight int, termWidth int, termHeight int, clientColor tcell.Color) *Box {
	return &Box{
		simWidth:    simWidth,
		simHeight:   simHeight,
		termWidth:   termWidth,
		termHeight:  termHeight,
		boxX:        0,
		boxY:        0,
		clientColor: clientColor,
		screen:      screen,
	}
}

// RenderBox draws the box on the screen.
// It also updates the boxX, boxY and keeps track of the new termWidth, termHeight
// This has to be called before SetContent() in the loop
func (b *Box) RenderBox(screen tcell.Screen, termWidth int, termHeight int) {
	// Calculate centered box position
	b.udpateBoxXY(termWidth, termHeight)
	b.termWidth = termWidth
	b.termHeight = termHeight

	// Draw box border
	borderStyle := tcell.StyleDefault.Foreground(b.clientColor)
	// Top and bottom borders
	for x := b.boxX; x < b.boxX+b.simWidth && x < termWidth; x++ {
		if b.boxY > 0 {
			screen.SetContent(x, b.boxY-1, '─', nil, borderStyle)
		}
		if b.boxY+b.simHeight < termHeight {
			screen.SetContent(x, b.boxY+b.simHeight, '─', nil, borderStyle)
		}
	}
	// Left and right borders
	for y := b.boxY; y < b.boxY+b.simHeight && y < termHeight; y++ {
		if b.boxX > 0 {
			screen.SetContent(b.boxX-1, y, '│', nil, borderStyle)
		}
		if b.boxX+b.simWidth < termWidth {
			screen.SetContent(b.boxX+b.simWidth, y, '│', nil, borderStyle)
		}
	}
	// Corners
	if b.boxX > 0 && b.boxY > 0 {
		screen.SetContent(b.boxX-1, b.boxY-1, '┌', nil, borderStyle)
	}
	if b.boxX+b.simWidth < termWidth && b.boxY > 0 {
		screen.SetContent(b.boxX+b.simWidth, b.boxY-1, '┐', nil, borderStyle)
	}
	if b.boxX > 0 && b.boxY+b.simHeight < termHeight {
		screen.SetContent(b.boxX-1, b.boxY+b.simHeight, '└', nil, borderStyle)
	}
	if b.boxX+b.simWidth < termWidth && b.boxY+b.simHeight < termHeight {
		screen.SetContent(b.boxX+b.simWidth, b.boxY+b.simHeight, '┘', nil, borderStyle)
	}
}

func (b *Box) InBox(x int, y int) bool {
	if x >= b.boxX && x < b.boxX+b.simWidth && y >= b.boxY && y < b.boxY+b.simHeight {
		return true
	}
	return false
}

func (b *Box) ToSimCoordinates(x int, y int) (float64, float64) {
	simX := float64((x - b.boxX) * 2)
	simY := float64((y - b.boxY) * 2)
	return simX, simY
}

func (b *Box) udpateBoxXY(termWidth int, termHeight int) {
	b.boxX = (termWidth - b.simWidth) / 2
	b.boxY = (termHeight - b.simHeight) / 2
	b.termWidth = termWidth
	b.termHeight = termHeight
	if b.boxX < 0 {
		b.boxX = 0
	}
	if b.boxY < 0 {
		b.boxY = 0
	}
}

func (b *Box) SetContent(x, y int, mainc rune, combc []rune, style tcell.Style) {
	if x >= 0 && x < b.simWidth && y >= 0 && y < b.simHeight {
		realX := x + b.boxX
		realY := y + b.boxY
		// Check bounds within the real screen view
		if realX >= 0 && realX < b.termWidth && realY >= 0 && realY < b.termHeight {
			b.screen.SetContent(realX, realY, mainc, combc, style)
		}
	}

}
