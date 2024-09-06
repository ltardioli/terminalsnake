package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	for {
		InitScreen()
		InitGameState()
		InitUserInput()

		for !isGameOver {
			ProcessInputs()
			UpdateState()
			DrawState()

			time.Sleep(75 * time.Millisecond)
		}

		DrawGameOver()

		// Wait for the user input after the game is over
		for isGameOver && !restart {
			ProcessInputs()
			time.Sleep(75 * time.Millisecond)
		}
		// Clean resources
		screen.Fini()

		if !restart {
			break
		}
	}
	screen.Fini()
}

func DrawState() {
	if isGamePaused {
		return
	}

	ClearScreen()
	PrintString(0, 0, debugLog)
	DrawGameFrame()

	DrawSnake()
	DrawApple()
	DrawScore()

	screen.Show()
}

func ClearScreen() {
	for _, p := range pointsToClear {
		DrawInsideGameFrame(p.row, p.col, 1, 1, ' ')
	}
	pointsToClear = []*Point{}
}

func UpdateState() {
	if isGamePaused {
		return
	}

	UpdateSnake()
	UpdateApples()
}

func UpdateApples() {
	// Remove eaten apples and calculate score. +1 normal apples +2 special apples
	// Put the elements that we want to keep at the beginning of the slice and then cut the slice
	n := 0
	for _, apple := range apples {
		if !IsAppleInsideSnake(apple) {
			apples[n] = apple
			n++
		} else {
			if apple.isSpecial && !apple.isTimed {
				score += 2
			} else if apple.isSpecial && apple.isTimed {
				score += 5
			} else {
				score++
			}
		}
	}
	apples = apples[:n]

	// Remove timed out apples
	for i := len(apples) - 1; i >= 0; i-- {
		if apples[i].isTimed && time.Since(apples[i].now) > AppleTimeout*time.Second {
			pointsToClear = append(pointsToClear, apples[i].point)
			apples = append(apples[:i], apples[i+1:]...)
		}
	}

	// Check how many simple apples needs to be generated and generate them
	count := 0
	for _, apple := range apples {
		if !apple.isSpecial {
			count++
		}
	}
	for i := 0; i < simultaneousApples-count; i++ {
		apples = append(apples, GenerateApple(AppleSymbol, false, false))
	}

	// Generate special apples
	specialApple := GenerateSpecialApple(SpecialAppleSymbol, false, SpecialAppleChance)
	if specialApple != nil {
		apples = append(apples, specialApple)
	}

	// Generate special timed apples
	specialTimedApple := GenerateSpecialApple(SpecialAppleSymbol, true, SpecialTimedAppleChance)
	if specialTimedApple != nil {
		apples = append(apples, specialTimedApple)
	}
}

func GenerateApple(symbol rune, isSpecial, isTimed bool) *Apple {
	apple := NewApple(NewPoint(rand.Intn(GameFrameHigh), rand.Intn(GameFrameWidth)), symbol, isSpecial, isTimed)
	for IsAppleInsideSnake(apple) {
		apple = NewApple(NewPoint(rand.Intn(GameFrameHigh), rand.Intn(GameFrameWidth)), symbol, isSpecial, isTimed)
	}
	return apple
}

func GenerateSpecialApple(symbol rune, timed bool, chance int) *Apple {
	luck := rand.Intn(chance)
	if luck == 1 { // At each frame drawed, generate a random number from 0 to SpecialAppleChance. If the generated number is 1 then generate a special apple. Chances by frame will be 1/SpecialAppleChance
		return GenerateApple(symbol, true, timed)
	}
	return nil
}

func IsAppleInsideSnake(apple *Apple) bool {
	for _, p := range snake.parts {
		if p.row == apple.point.row && p.col == apple.point.col {
			return true
		}
	}
	return false
}

func UpdateSnake() {
	// The snake will moving by adding a new head in the new direction and removing the tail
	head := GetSnakeHead()
	snake.parts = append(snake.parts, &Point{row: head.row + snake.velRow, col: head.col + snake.velCol})

	// Growing the snake if it eats an apple. Just do not delete the tail when updates the moviment of the snake
	appleEaten := false
	for _, apple := range apples {
		if IsAppleInsideSnake(apple) {
			appleEaten = true
		}
	}
	if !appleEaten {
		snake.parts = snake.parts[1:]
	}

	// Check snake colisions
	if IsSnakeHittingWall() || IsSnakeEatingItself() {
		isGameOver = true
	}
}

func IsSnakeEatingItself() bool {
	head := GetSnakeHead()
	for i := 0; i < len(snake.parts)-1; i++ {
		if head.col == snake.parts[i].col && head.row == snake.parts[i].row {
			return true
		}
	}

	return false
}

func IsSnakeHittingWall() bool {
	head := GetSnakeHead()
	return head.row < 0 || head.row >= GameFrameHigh || head.col < 0 || head.col >= GameFrameWidth
}

func GetSnakeHead() *Point {
	return snake.parts[len(snake.parts)-1]
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.HideCursor()
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorNone)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	snake = &Snake{
		parts: []*Point{
			{row: 5, col: 3},
			{row: 6, col: 3},
			{row: 7, col: 3},
			{row: 8, col: 3},
			{row: 9, col: 3},
		},
		velRow: 1,
		velCol: 0,
		symbol: SnakeSymbol,
	}

	simultaneousApples = 1 // The idea is to add more apples at the same time in the future and it will increase and decrease dynamically
	apples = []*Apple{
		{
			point:     &Point{row: 10, col: 10},
			symbol:    AppleSymbol,
			isSpecial: false,
		},
	}
	score = 0
	restart = false
	isGameOver = false
}

func DrawSnake() {
	for _, p := range snake.parts {
		DrawInsideGameFrame(p.row, p.col, 1, 1, snake.symbol)
		pointsToClear = append(pointsToClear, p)
	}
}

func DrawApple() {
	for _, apple := range apples {
		color := Red
		if apple.isSpecial && !apple.isTimed {
			color = Green
		} else if apple.isSpecial && apple.isTimed {
			color = Yellow
		}
		DrawInsideGameFrame(apple.point.row, apple.point.col, 1, 1, apple.symbol, color)
		pointsToClear = append(pointsToClear, apple.point)
	}
}

func DrawScore() {
	row, col := GetGameFrameTopLeft()
	PrintString(row-2, col, fmt.Sprintf("Score: %d", score))
}

func DrawGameOver() {
	screenWidth, screenHeight := screen.Size()
	PrintStringCentered(screenHeight/2, screenWidth/2, "Game Over!")
	PrintStringCentered(screenHeight/2+1, screenWidth/2+1, fmt.Sprint("Your score is: ", score))
	PrintStringCentered(screenHeight/2+2, screenWidth/2+2, "Press Esc or 'q' to leave")
	PrintStringCentered(screenHeight/2+3, screenWidth/2+3, "Press Enter to try again!")
	screen.Show()
}

func DrawInsideGameFrame(row, col, width, height int, ch rune, color ...Color) {
	rowOffset, colOffset := GetGameFrameTopLeft()
	DrawFilledRect(row+rowOffset, col+colOffset, width, height, ch, color...)
}

func DrawFilledRect(row, col, width, height int, ch rune, color ...Color) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			var style tcell.Style
			if color == nil {
				style = tcell.StyleDefault
			} else {
				style = GetColor(color[0])
			}
			screen.SetContent(col+c, row+r, ch, nil, style)
		}
	}
}

func DrawUnfilledRect(row, col, width, height int, ch rune, color ...Color) {
	var style tcell.Style
	if color == nil {
		style = tcell.StyleDefault
	} else {
		style = GetColor(color[0])
	}

	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, style)
		screen.SetContent(col+c, row+height-1, ch, nil, style)
	}

	for r := 0; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, style)
		screen.SetContent(col+width-1, row+r, ch, nil, style)
	}
}

func DrawGameFrame() {
	gameFrameTopLeftRow, gameFrameTopLeftCol := GetGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, height := GameFrameWidth+2, GameFrameHigh+2

	DrawUnfilledRect(row, col, width, height, GameFrameSymbol)
	//DrawUnfilledRect(row+1, col+1, GameFrameWidth, GameFrameHigh, '*')
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func GetGameFrameTopLeft() (int, int) {
	screnWidth, screenHeight := screen.Size()
	return (screenHeight - GameFrameHigh) / 2, (screnWidth - GameFrameWidth) / 2
}

func InitUserInput() {
	go func() {
		for {
			switch ev := screen.PollEvent().(type) { // Block waiting for the event
			case *tcell.EventKey:
				//debugLog = ev.Name()
				key := ev.Name()
				mu.Lock()
				// Avoid subsequently equals inputs if they are the same. Avoid buffering in case the user holds a button.
				if len(inputs) == 0 || (len(inputs) > 0 && inputs[len(inputs)-1] != key) {
					inputs = append(inputs, key)
				}
				mu.Unlock()
			}
		}
	}()
}

func ProcessInputs() {
	mu.Lock()
	if len(inputs) > 0 {
		HandleUserInput(inputs[0])
		inputs = inputs[1:]
	}
	mu.Unlock()
}

func HandleUserInput(key string) {
	if key == "Rune[q]" || key == " " || key == "Esc" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if (key == "Right" || key == "Rune[d]") && snake.velCol != -1 {
		snake.velCol = 1
		snake.velRow = 0
	} else if (key == "Up" || key == "Rune[w]") && snake.velRow != 1 {
		snake.velCol = 0
		snake.velRow = -1
	} else if (key == "Down" || key == "Rune[s]") && snake.velRow != -1 {
		snake.velCol = 0
		snake.velRow = 1
	} else if (key == "Left" || key == "Rune[a]") && snake.velCol != 1 {
		snake.velCol = -1
		snake.velRow = 0
	} else if key == "Enter" && isGameOver {
		restart = true
	}

}
