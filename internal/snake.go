package internal

import "sync"

type SnakeCell struct {
	positionX     int32
	positionY     int32
	prevPositionX int32
	prevPositionY int32
	symbol        int32
}

type Coordinates struct {
	positionX int32
	positionY int32
}

type SnakeStats struct {
	countOfEatBoosters int32
	size               int32
}

type Snake struct {
	sync.Mutex
	cells []*SnakeCell
	head  *SnakeCell
	SnakeStats
}

func CreateSnake() *Snake {
	snake := new(Snake)
	snake.Init()

	return snake
}

func (s *Snake) Init() {
	s.size = 0
}

func (s *Snake) CreateHead(symbol int32, x int32, y int32) *SnakeCell {
	cell := new(SnakeCell)

	cell.symbol = symbol
	cell.positionX = x
	cell.positionY = y

	s.head = cell

	return cell
}

func (s *Snake) IncrementAteBoosters() {
	s.countOfEatBoosters++
}

func (s *Snake) GetCountOfEatBoosters() int32 {
	s.Lock()
	defer s.Unlock()

	return s.countOfEatBoosters
}

func (s *Snake) IncrementSize() {
	s.Lock()
	defer s.Unlock()

	s.size++
}

func (s *Snake) AddCell(symbol int32) *SnakeCell {
	s.Lock()
	defer s.Unlock()

	cell := new(SnakeCell)

	cell.symbol = symbol

	s.cells = append(s.cells, cell)
	s.size++

	return cell
}

func (s *Snake) GetCells() []*SnakeCell {
	s.Lock()
	defer s.Unlock()

	return s.cells
}

func (s *Snake) GetFirstCell() *SnakeCell {
	s.Lock()
	defer s.Unlock()

	if len(s.cells) > 0 {
		return s.cells[0]
	}

	return nil
}

func (s *Snake) GetHeadCell() *SnakeCell {
	s.Lock()
	defer s.Unlock()

	return s.head
}

func (s *Snake) isHeadOnBody() bool {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.cells {
		if s.head.positionX == v.positionX && s.head.positionY == v.positionY {
			return true
		}
	}

	return false
}

func (s *Snake) GetSize() int32 {
	s.Lock()
	defer s.Unlock()

	return s.size
}

func (s *Snake) CalculatePositions(cell SnakeCell, direction int32) SnakeCell {
	s.Lock()
	defer s.Unlock()

	switch direction {
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

	return cell
}

func (s *Snake) DeleteLastCells(c int) []Coordinates {
	s.Lock()
	defer s.Unlock()

	index := len(s.cells) - c

	if index < 0 {
		index = 0
	}

	var tmp = *new([]*SnakeCell)
	var coordinates = *new([]Coordinates)

	j := 0

	for i, cell := range s.cells {
		if i >= index {
			coords := new(Coordinates)
			coords.positionX = cell.positionX
			coords.positionY = cell.positionY

			coordinates = append(coordinates, *coords)

			j++
			s.size--

			continue
		}

		tmp = append(tmp, cell)
	}

	s.cells = tmp

	return coordinates
}
