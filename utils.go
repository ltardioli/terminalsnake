package main

import "github.com/gdamore/tcell/v2"

func GetColor(color Color) tcell.Style {
	switch color {
	case White:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	case Black:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	case Red:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	case Blue:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	case Green:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	case Yellow:
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
	default:
		return tcell.StyleDefault
	}
}
