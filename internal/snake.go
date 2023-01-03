package internal

type SnakeCell struct {
	positionX     int32
	positionY     int32
	prevPositionX int32
	prevPositionY int32
	symbol        int32
}

type Snake struct {
	size          int32
	cells         []*SnakeCell
	head          *SnakeCell
	isContainBody bool
}

func (s *Snake) Init() {
	s.size = 0
}

func (s *Snake) AddCell(symbol int32) *SnakeCell {
	cell := new(SnakeCell)

	cell.symbol = symbol

	s.cells = append(s.cells, cell)
	s.size++
	s.isContainBody = true

	return cell
}

func (s *Snake) GetFirstCell() *SnakeCell {
	if s.isContainBody {
		return s.cells[0]
	}

	return nil
}

func (s *Snake) CreateHead(symbol int32, x int32, y int32) *SnakeCell {
	cell := new(SnakeCell)

	cell.symbol = symbol
	cell.positionX = x
	cell.positionY = y

	s.head = cell

	return cell
}

func (s *Snake) GetHeadCell() *SnakeCell {
	return s.head
}

func (s *Snake) isHeadOnBody() bool {
	for _, v := range s.cells {
		if s.head.positionX == v.positionX && s.head.positionY == v.positionY {
			return true
		}
	}

	return false
}

func (s *Snake) GetSize() int32 {
	return s.size
}

func (s *Snake) CalculatePositions(cell SnakeCell, direction int32) SnakeCell {
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
