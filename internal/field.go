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
	"os"
)

const lineSize int32 = 10
const isAllowCrossBoarding = true

const directions int32 = 0x04

const fieldSymbol = '.'
const headSymbol = 0x40
const bodySymbol = 0x23

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
	f.isAllowCrossBoarding = isAllowCrossBoarding
	f.buffer = new([lineSize][lineSize]string)
	f.logger = l

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
}

func (f *Field) SetSnakeHead() {
	min := int32(2)

	posX := rand.Int31n(lineSize-min) + 1
	posY := rand.Int31n(lineSize-min) + 1

	f.direction = rand.Int31n(directions)
	cell := f.snake.AddCell(headSymbol, posX, posY, true)

	f.buffer[cell.positionX][cell.positionY] = string(cell.symbol)
}

func (f *Field) MoveHead() {
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

	if f.snake.isHeadOnBody(cell.positionX, cell.positionY) {
		panic("you was failed")
	}

	file, _ := os.OpenFile("./test.log", os.O_APPEND|os.O_WRONLY, os.ModePerm)

	file.WriteString(fmt.Sprint(cell.positionX, cell.positionY, f.boosterPositionX, f.boosterPositionY, "\n"))

	if cell.positionX == f.boosterPositionX && cell.positionY == f.boosterPositionY {
		file.WriteString(fmt.Sprint(cell.prevPositionX, cell.prevPositionY, "\n"))

		f.snake.AddCell(bodySymbol, cell.prevPositionX, cell.prevPositionY, false)

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

	f.snake.cells[0] = cell
	f.buffer[cell.positionX][cell.positionY] = string(cell.symbol)

	file.Close()
}

func (f *Field) FillBody() {
	headCell := f.snake.GetHeadCell()

	prevPositionX := headCell.prevPositionX
	prevPositionY := headCell.prevPositionY

	file, _ := os.OpenFile("./test.log", os.O_APPEND|os.O_WRONLY, os.ModePerm)

	defer file.Close()

	for i, v := range f.snake.cells {
		file.WriteString(fmt.Sprint("fillbody: ", v.isHead, v.prevPositionX, v.prevPositionY, v.positionX, v.positionY, "\n"))

		if v.isHead {
			continue
		}

		v.prevPositionX = v.positionX
		v.prevPositionY = v.positionY

		v.positionX = prevPositionX
		v.positionY = prevPositionY

		f.snake.cells[i] = v

		f.buffer[v.prevPositionX][v.prevPositionY] = string(fieldSymbol)
		f.buffer[v.positionX][v.positionY] = string(v.symbol)

		prevPositionX = v.prevPositionX
		prevPositionY = v.prevPositionY
	}
}

func (f *Field) OnStep() error {
	f.MoveHead()
	f.FillBody()

	return nil

	for i, v := range f.snake.cells {
		v.prevPositionX = v.positionX
		v.prevPositionY = v.positionY

		f.buffer[v.positionX][v.positionY] = string(fieldSymbol)

		if v.isHead {
			switch f.direction {
			case directionLeft:
				v.positionY--
				break
			case directionUp:
				v.positionX--
				break
			case directionRight:
				v.positionY++
				break
			case directionDown:
				v.positionX++
				break
			}

			if f.snake.isContainBody {
				nextCell := f.snake.cells[i+1]

				if v.positionX == nextCell.positionX && v.positionY == nextCell.positionY {
					v.positionX = v.prevPositionX
					v.positionY = v.prevPositionY

					f.direction = f.prevDirection

					switch f.direction {
					case directionLeft:
						v.positionY--
						break
					case directionUp:
						v.positionX--
						break
					case directionRight:
						v.positionY++
						break
					case directionDown:
						v.positionX++
						break
					}
				}
			}

			if f.snake.isHeadOnBody(v.positionX, v.positionY) {
				return errors.New("you was failed")
			}

			if v.positionX == f.boosterPositionX && v.positionY == f.boosterPositionY {
				f.snake.AddCell(bodySymbol, v.prevPositionX, v.prevPositionY, false)

				f.isBoosterSpawned = false
			}
		} else {
			vOld := f.snake.cells[i-1]

			v.positionX = vOld.prevPositionX
			v.positionY = vOld.prevPositionY
		}

		if v.positionX == xBorderLower && f.direction == directionUp {
			if !isAllowCrossBoarding {
				return errors.New("you was failed")
			}

			v.positionX = xBorderHigherPossible
		}

		if v.positionX == xBorderHigher && f.direction == directionDown {
			if !isAllowCrossBoarding {
				return errors.New("you was failed")
			}

			v.positionX = xBorderLowerPossible
		}

		if v.positionY == yBorderLower && f.direction == directionLeft {
			if !isAllowCrossBoarding {
				return errors.New("you was failed")
			}

			v.positionY = yBorderHigherPossible
		}

		if v.positionY == yBorderHigher && f.direction == directionRight {
			if !isAllowCrossBoarding {
				return errors.New("you was failed")
			}

			v.positionY = yBorderLowerPossible
		}

		f.snake.cells[i] = v
		f.buffer[v.positionX][v.positionY] = string(v.symbol)
	}

	return nil
}

func (f *Field) ChangeDirection(d int32) {
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

	f.buffer[spawnPointX][spawnPointY] = "%"
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

func (f *Field) GetSnake() *Snake {
	return f.snake
}
