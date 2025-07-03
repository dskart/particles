package app

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

func colorName(color tcell.Color) string {
	switch color {
	case tcell.ColorRed:
		return "Red"
	case tcell.ColorGreen:
		return "Green"
	case tcell.ColorBlue:
		return "Blue"
	case tcell.ColorYellow:
		return "Yellow"
	case tcell.ColorFuchsia:
		return "Magenta"
	case tcell.ColorAqua:
		return "Cyan"
	case tcell.ColorOrange:
		return "Orange"
	case tcell.ColorPurple:
		return "Purple"
	default:
		return "Unknown"
	}
}

var clientColors = []tcell.Color{
	tcell.ColorRed,
	tcell.ColorGreen,
	tcell.ColorBlue,
	tcell.ColorYellow,
	tcell.ColorFuchsia,
	tcell.ColorAqua,
	tcell.ColorOrange,
	tcell.ColorPurple,
}

func getRandColor() tcell.Color {
	return clientColors[rand.Intn(len(clientColors))]
}
