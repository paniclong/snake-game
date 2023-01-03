package internal

/*
#include <curses.h>
#include <stdio.h>
#cgo LDFLAGS: -lcurses
*/
import "C"
import (
	"errors"
	"fmt"
	"math/rand"
)

const lineSize int32 = 10
const isAllowCrossBoarding = true

const directions int32 = 0x04

const fieldSymbol = '.'
const headSymbol = 0x40
const bodySymbol = 0x23
const boosterSymbol = 0x25

const (
	directionLeft  = 0
	directionUp    = 1
	directionRight = 2
	directionDown  = 3
)

const (
	keyBackspace = 263
	keyLeft      = 260
	keyUp        = 259
	keyRight     = 261
	keyDown      = 258
)

// Coordinates borders (x, y). Needed when flag isAllowCrossBoarding = true
// Allows snakes to move to the opposite border
const (
	xBorderLower          = 0x0
	yBorderLower          = 0x0
	xBorderLowerPossible  = xBorderLower + 1
	yBorderLowerPossible  = yBorderLower + 1
	xBorderHigher         = lineSize - 1
	yBorderHigher         = lineSize - 1
	xBorderHigherPossible = xBorderHigher - 1
	yBorderHigherPossible = yBorderHigher - 1
)

type Field struct {
	snake                *Snake
	buffer               *[lineSize][lineSize]string
	direction            int32 // Default direction (left - 0, up - 1, right - 2, down - 3)
	prevDirection        int32
	isActive             bool
	isBoosterSpawned     bool
	boosterPositionX     int32
	boosterPositionY     int32
	isAllowCrossBoarding bool
	logger               *Logger
}

func (f *Field) Init(s *Snake, l *Logger) {
	f.isActive = true
	f.snake = s
	f.logger = l
	f.isAllowCrossBoarding = isAllowCrossBoarding
	f.buffer = new([lineSize][lineSize]string)

	for i, line := range f.buffer {
		for j := range line {
			f.buffer[i][j] = string(fieldSymbol)
		}
	}

	f.SetSnakeHead()
}

func (f *Field) Print() {
	for i, line := range f.buffer {
		i := int32(i)

		if i == 0 || i == (lineSize-1) {
			for j := int32(0); j < lineSize; j++ {
				C.addstr(C.CString(" "))
				C.addstr(C.CString("-"))
				C.addstr(C.CString(" "))
			}

			C.addstr(C.CString("\n"))

			continue
		}

		for j, v := range line {
			j := int32(j)

			if j == 0 || j == (lineSize-1) {
				C.addstr(C.CString(" "))
				C.addstr(C.CString("|"))
				C.addstr(C.CString(" "))

				continue
			}

			C.addstr(C.CString(" "))
			C.addstr(C.CString(v))
			C.addstr(C.CString(" "))
		}

		C.addstr(C.CString("\n"))
	}

	C.addstr(C.CString("Press backspace for exit\n"))
}

func (f *Field) SetSnakeHead() {
	min := int32(2)

	posX := rand.Int31n(lineSize-min) + 1
	posY := rand.Int31n(lineSize-min) + 1

	f.direction = rand.Int31n(directions)
	cell := f.snake.CreateHead(headSymbol, posX, posY)

	f.buffer[cell.positionX][cell.positionY] = string(cell.symbol)
}

func (f *Field) MoveHead() error {
	cell := f.snake.GetHeadCell()

	// Current coordinates writing to prev properties
	cell.prevPositionX = cell.positionX
	cell.prevPositionY = cell.positionY

	f.buffer[cell.positionX][cell.positionY] = string(fieldSymbol)

	switch f.direction {
	case directionLeft:
		cell.positionY--
		break
	case directionUp:
		cell.positionX--
		break
	case directionRight:
		cell.positionY++
		break
	case directionDown:
		cell.positionX++
		break
	}

	if f.snake.isHeadOnBody() {
		return errors.New("you ate yourself")
	}

	f.logger.WriteString(
		fmt.Sprintf(
			"Calculated positions for head: currentX [%d], currentY [%d], prevX [%d], prevY [%d]",
			cell.positionX, cell.positionY, cell.prevPositionX, cell.prevPositionY,
		),
	)

	if cell.positionX == f.boosterPositionX && cell.positionY == f.boosterPositionY {
		f.logger.WriteString(fmt.Sprint("Booster eats"))

		f.snake.AddCell(bodySymbol)

		f.isBoosterSpawned = false
	}

	if cell.positionX == xBorderLower && f.direction == directionUp {
		cell.positionX = xBorderHigherPossible
	}

	if cell.positionX == xBorderHigher && f.direction == directionDown {
		cell.positionX = xBorderLowerPossible
	}

	if cell.positionY == yBorderLower && f.direction == directionLeft {
		cell.positionY = yBorderHigherPossible
	}

	if cell.positionY == yBorderHigher && f.direction == directionRight {
		cell.positionY = yBorderLowerPossible
	}

	f.buffer[cell.positionX][cell.positionY] = string(cell.symbol)

	return nil
}

func (f *Field) FillBody() {
	headCell := f.snake.GetHeadCell()

	prevPositionX := headCell.prevPositionX
	prevPositionY := headCell.prevPositionY

	for i, v := range f.snake.cells {
		f.logger.WriteString(
			fmt.Sprintf(
				"Fill body: i [%d], currentX [%d], currentY [%d], prevX [%d], prevY [%d], prevHeadX [%d], prevHeadY [%d]",
				i, v.positionX, v.positionY, v.prevPositionX, v.prevPositionY, prevPositionX, prevPositionY,
			),
		)

		v.prevPositionX = v.positionX
		v.prevPositionY = v.positionY

		v.positionX = prevPositionX
		v.positionY = prevPositionY

		f.buffer[v.prevPositionX][v.prevPositionY] = string(fieldSymbol)
		f.buffer[v.positionX][v.positionY] = string(v.symbol)

		prevPositionX = v.prevPositionX
		prevPositionY = v.prevPositionY

		f.logger.WriteString(
			fmt.Sprintf(
				"Fill body end: i [%d], currentX [%d], currentY [%d], prevX [%d], prevY [%d], prevHeadX [%d], prevHeadY [%d]",
				i, v.positionX, v.positionY, v.prevPositionX, v.prevPositionY, prevPositionX, prevPositionY,
			),
		)
	}
}

func (f *Field) OnStep() error {
	err := f.MoveHead()
	if err != nil {
		return err
	}

	f.FillBody()

	return nil
}

func (f *Field) ChangeDirection(d int32) {
	headCell := *f.snake.GetHeadCell()
	headCell = f.snake.CalculatePositions(headCell, d)

	cell := f.snake.GetFirstCell()

	// If direction was changed to a cell, that contains a body cell, this action is rejected
	if cell != nil && headCell.positionX == cell.positionX && headCell.positionY == cell.positionY {
		return
	}

	f.prevDirection = f.direction
	f.direction = d
}

func (f *Field) SpawnBooster() {
	if f.isBoosterSpawned {
		return
	}

	min := int32(2)

	var spawnPointX int32 = 0
	var spawnPointY int32 = 0

	for {
	repeat:
		spawnPointX = rand.Int31n(lineSize-min) + 1
		spawnPointY = rand.Int31n(lineSize-min) + 1

		for _, v := range f.snake.cells {
			if v.positionX == spawnPointX && v.positionY == spawnPointY {
				goto repeat
			}
		}

		break
	}

	f.boosterPositionX = spawnPointX
	f.boosterPositionY = spawnPointY
	f.isBoosterSpawned = true

	f.buffer[spawnPointX][spawnPointY] = string(rune(boosterSymbol))
}

func (f *Field) ChangeDirectionByKey() {
	for {
		key := C.getch()

		var d int32 = -1

		switch key {
		case keyBackspace:
			f.isActive = false
		case keyLeft:
			d = directionLeft
			break
		case keyUp:
			d = directionUp
			break
		case keyRight:
			d = directionRight
			break
		case keyDown:
			d = directionDown
			break
		}

		if d != -1 {
			f.ChangeDirection(d)
		}
	}
}

func (f *Field) IsActive() bool {
	return f.isActive
}

func (f *Field) SetIsActive(isActive bool) {
	f.isActive = isActive
}

func (f *Field) GetSnake() *Snake {
	return f.snake
}
