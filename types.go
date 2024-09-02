package main

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
}

func NewApple(point *Point, symbol rune, isSpecial bool) *Apple {
	return &Apple{
		point:     point,
		symbol:    symbol,
		isSpecial: isSpecial,
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
