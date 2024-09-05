package main

import "time"

type Point struct {
	row, col int
}

func NewPoint(row, col int) *Point {
	return &Point{
		row: row,
		col: col,
	}
}

type Snake struct {
	parts          []*Point
	velRow, velCol int
	symbol         rune
}

type Apple struct {
	point     *Point
	symbol    rune
	isSpecial bool
	isTimed   bool
	now       time.Time
}

func NewApple(point *Point, symbol rune, isSpecial, isTimed bool) *Apple {
	return &Apple{
		point:     point,
		symbol:    symbol,
		isSpecial: isSpecial,
		isTimed:   isTimed,
		now:       time.Now(),
	}
}

type Color int

const (
	White Color = iota
	Black
	Blue
	Red
	Green
	Yellow
)
