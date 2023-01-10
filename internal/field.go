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
	"time"
)

// Size of field (lineSize X lineSize)
// For comforting playing lineSize should be greater or equals 15
const lineSize int32 = 20
const isAllowCrossBoarding = false
const limitBoosters = 1
const limitEnemies = 1

const startSnakeLength = 3
const directions int32 = 0x04

const fieldSymbol = 0x2e
const xBorderSymbol = 0x2d
const yBorderSymbol = 0x7c

const HeadSymbol = 0x40
const BodySymbol = 0x23
const BoosterSymbol = 0x25
const EnemyOneShotSymbol = 0x21
const EnemySymbol = 0x24

const spaceSymbol = 0x20

const (
	directionLeft  = iota
	directionUp    = iota
	directionRight = iota
	directionDown  = iota
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
	buffer               *[lineSize][lineSize]int32
	direction            int32 // Default direction (left - 0, up - 1, right - 2, down - 3)
	prevDirection        int32
	isActive             bool
	IsFirstStart         bool
	isAllowCrossBoarding bool
	logger               *Logger
	boosters             []*Booster
	enemies              []*Enemy
}

func (f *Field) Init(s *Snake, l *Logger) {
	f.isActive = true
	f.IsFirstStart = true
	f.snake = s
	f.logger = l
	f.isAllowCrossBoarding = isAllowCrossBoarding
	f.buffer = new([lineSize][lineSize]int32)

	for x, line := range f.buffer {
		x := int32(x)

		for y := range line {
			y := int32(y)

			if x == 0 || x == (lineSize-1) {
				f.buffer[x][y] = xBorderSymbol

				continue
			}

			if y == 0 || y == (lineSize-1) {
				f.buffer[x][y] = yBorderSymbol

				continue
			}

			f.buffer[x][y] = fieldSymbol
		}
	}

	f.InitSnake()
}

func (f *Field) Print() {
	var s string

	for _, line := range f.buffer {
		for _, v := range line {

			s += string(spaceSymbol) + string(v) + string(spaceSymbol)
		}

		s += "\n"
	}

	s += fmt.Sprintf("Symbols: \n 1. %s - snake head \n 2. %s - snake body \n 3. %s - boosters \n "+
		"4. %s - enemy, one shot \n 5. %s - enemy, no one shot \n",
		string(HeadSymbol),
		string(BodySymbol),
		string(BoosterSymbol),
		string(EnemyOneShotSymbol),
		string(EnemySymbol),
	)

	s += "\n"
	s += fmt.Sprintf("Snake size: %d\n", f.snake.size+1)
	s += fmt.Sprintf("Number of boosters eaten: %d\n", f.snake.countOfEatBoosters)
	for i, v := range f.enemies {
		s += fmt.Sprintf("Enemy [%d], isOneShot: [%t], cells: [%d]\n", i, v.isOneShot, v.countOfCells)
	}
	s += "\n"
	s += "Press backspace for exit\n"

	C.addstr(C.CString(s))

	f.logger.WriteString(s)
}

func (f *Field) InitSnake() {
	var respawnOffset int32 = 3

	posX := RandomInt32MinMaxN(xBorderLower+(xBorderHigher/2), (xBorderHigher/2)+respawnOffset)
	posY := RandomInt32MinMaxN(yBorderLower+(yBorderHigher/2), (xBorderHigher/2)+respawnOffset)

	head := f.snake.CreateHead(HeadSymbol, posX, posY)

	f.buffer[head.positionX][head.positionY] = head.symbol

	f.logger.WriteString(fmt.Sprintf("Spawn head %v", head))

	l := 0

	for {
		if l >= startSnakeLength {
			break
		}

		posX--

		bodyCell := f.snake.AddCell(BodySymbol)

		bodyCell.positionX = posX
		bodyCell.positionY = posY

		l++

		f.buffer[bodyCell.positionX][bodyCell.positionY] = bodyCell.symbol

		f.logger.WriteString(fmt.Sprintf("Spawn body %v", bodyCell))
	}

	d := rand.Int31n(directions)

	for {
		err := f.ChangeDirection(d)

		if err == nil {
			break
		}

		d--
	}
}

func (f *Field) MoveHead() error {
	cell := f.snake.GetHeadCell()

	// Current coordinates writing to prev properties
	cell.prevPositionX = cell.positionX
	cell.prevPositionY = cell.positionY

	f.buffer[cell.positionX][cell.positionY] = fieldSymbol

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
		true,
	)

	for i, booster := range f.boosters {
		if cell.positionX == booster.positionX && cell.positionY == booster.positionY {
			f.logger.WriteString(fmt.Sprintf("Booster [%d] eats", i))

			f.snake.countOfEatBoosters++
			f.snake.AddCell(BodySymbol)

			f.DeSpawnBooster(i)
		}
	}

	for i, enemy := range f.enemies {
		if cell.positionX == enemy.positionX && cell.positionY == enemy.positionY && enemy.isActive {
			countOfCells := enemy.countOfCells
			if enemy.isOneShot || len(f.snake.cells) < int(countOfCells) {
				f.logger.WriteString(fmt.Sprintf("Enemy [%d] killed you", i))

				return errors.New(fmt.Sprintf("Enemy [%d] killed you", i))
			} else {
				coordinates := f.snake.DeleteLastCells(int(countOfCells))

				f.logger.WriteString(fmt.Sprintf("Delete last cells: %v", coordinates))

				for _, v := range coordinates {
					f.buffer[v.positionX][v.positionY] = fieldSymbol
				}
			}
		}
	}

	if isAllowCrossBoarding {
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
	} else {
		if cell.positionX == xBorderLower ||
			cell.positionY == yBorderLower ||
			cell.positionX == xBorderHigher ||
			cell.positionY == yBorderHigher {
			f.logger.WriteString("Cross boarding is not allowed, you failed")

			return errors.New("borders")
		}
	}

	f.buffer[cell.positionX][cell.positionY] = cell.symbol

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
			true,
		)

		v.prevPositionX = v.positionX
		v.prevPositionY = v.positionY

		v.positionX = prevPositionX
		v.positionY = prevPositionY

		f.buffer[v.prevPositionX][v.prevPositionY] = fieldSymbol
		f.buffer[v.positionX][v.positionY] = v.symbol

		prevPositionX = v.prevPositionX
		prevPositionY = v.prevPositionY

		f.logger.WriteString(
			fmt.Sprintf(
				"Fill body end: i [%d], currentX [%d], currentY [%d], prevX [%d], prevY [%d], prevHeadX [%d], prevHeadY [%d]",
				i, v.positionX, v.positionY, v.prevPositionX, v.prevPositionY, prevPositionX, prevPositionY,
			),
			true,
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

func (f *Field) ChangeDirection(d int32) error {
	headCell := *f.snake.GetHeadCell()
	headCell = f.snake.CalculatePositions(headCell, d)

	cell := f.snake.GetFirstCell()

	// If direction was changed to a cell, that contains a body cell, this action is rejected
	if cell != nil && headCell.positionX == cell.positionX && headCell.positionY == cell.positionY {
		return errors.New("wrong direction")
	}

	f.prevDirection = f.direction
	f.direction = d

	f.logger.WriteString(fmt.Sprintf("Changed direction, now: [%d], previos: [%d]", f.direction, f.prevDirection))

	return nil
}

func (f *Field) ChangeDirectionByKey() {
	for {
		key := C.getch()

		f.logger.WriteString(fmt.Sprint("ChangeDirectionByKey start ----> ", key))

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

		f.logger.WriteString(fmt.Sprint("ChangeDirectionByKey ", d))

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

func (f *Field) ReCalcEnemies() {
	for i, v := range f.GetEnemies() {
		if v.GetSpawnTime()+10 < time.Now().Unix() {
			f.logger.WriteString(fmt.Sprintf("Enemy [%d] has in field is too long, despawn", i))

			f.DeSpawnEnemy(i)
		}
	}

	f.SpawnEnemy()
}

func (f *Field) ReCalcBoosters() {
	f.SpawnBooster()
}
