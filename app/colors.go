package app

import (
	"math/rand/v2"
	"sync"

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
	case tcell.ColorLime:
		return "Lime"
	case tcell.ColorMaroon:
		return "Maroon"
	case tcell.ColorNavy:
		return "Navy"
	case tcell.ColorOlive:
		return "Olive"
	case tcell.ColorSilver:
		return "Silver"
	case tcell.ColorTeal:
		return "Teal"
	case tcell.ColorWhite:
		return "White"
	case tcell.ColorGray:
		return "Gray"
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
	tcell.ColorLime,
	tcell.ColorMaroon,
	tcell.ColorNavy,
	tcell.ColorOlive,
	tcell.ColorSilver,
	tcell.ColorTeal,
	tcell.ColorWhite,
	tcell.ColorGray,
}

type ColorManager struct {
	mu               sync.Mutex
	availableIndexes []int
	assignedColors   map[tcell.Color]bool
}

var globalColorManager = &ColorManager{
	availableIndexes: make([]int, len(clientColors)),
	assignedColors:   make(map[tcell.Color]bool),
}

func init() {
	for i := range clientColors {
		globalColorManager.availableIndexes[i] = i
	}
}

func GetAvailableColor() (tcell.Color, bool) {
	globalColorManager.mu.Lock()
	defer globalColorManager.mu.Unlock()

	if len(globalColorManager.availableIndexes) == 0 {
		return tcell.ColorDefault, false
	}

	index := rand.IntN(len(globalColorManager.availableIndexes))
	colorIndex := globalColorManager.availableIndexes[index]
	color := clientColors[colorIndex]

	globalColorManager.availableIndexes = append(
		globalColorManager.availableIndexes[:index],
		globalColorManager.availableIndexes[index+1:]...,
	)
	globalColorManager.assignedColors[color] = true

	return color, true
}

func ReleaseColor(color tcell.Color) {
	globalColorManager.mu.Lock()
	defer globalColorManager.mu.Unlock()

	if globalColorManager.assignedColors[color] {
		delete(globalColorManager.assignedColors, color)
		for i, c := range clientColors {
			if c == color {
				globalColorManager.availableIndexes = append(globalColorManager.availableIndexes, i)
				break
			}
		}
	}
}

func getRandColor() tcell.Color {
	color, ok := GetAvailableColor()
	if !ok {
		return clientColors[rand.IntN(len(clientColors))]
	}
	return color
}
