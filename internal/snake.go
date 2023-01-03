package internal

type SnakeCell struct {
	positionX     int32
	positionY     int32
	prevPositionX int32
	prevPositionY int32
	isHead        bool
	symbol        int32
}

type Snake struct {
	head          int32
	size          int32
	cells         []SnakeCell
	isContainBody bool
}

func (s *Snake) Init() {
	s.head = headSymbol
	s.size = 0
}

func (s *Snake) AddCell(symbol int32, x int32, y int32, isHead bool) SnakeCell {
	cell := new(SnakeCell)

	cell.positionX = x
	cell.positionY = y
	cell.isHead = isHead
	cell.symbol = symbol

	if isHead != true {
		s.isContainBody = true
	}

	s.cells = append(s.cells, *cell)
	s.size++

	return *cell
}

func (s *Snake) GetHeadCell() SnakeCell {
	for _, v := range s.cells {
		if v.isHead {
			return v
		}
	}

	panic("Cannot get head cell")
}

func (s *Snake) isHeadOnBody(x int32, y int32) bool {
	for _, v := range s.cells {
		if v.isHead {
			continue
		}

		if x == v.positionX && y == v.positionY {
			return true
		}
	}

	return false
}

func (s *Snake) GetSize() int32 {
	return s.size
}
