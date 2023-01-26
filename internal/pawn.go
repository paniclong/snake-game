package internal

import (
	"fmt"
	"sync"
	"time"
)

type Pawn struct {
	sync.Mutex
	positionX int32
	positionY int32
	symbol    int32
	spawnTime int64
}

type Booster struct {
	Pawn
	isActive bool
}

type Enemy struct {
	Pawn
	isActive     bool
	isOneShot    bool
	countOfCells int32 // if isOneShot is false we are determine how many cells delete from snake
}

func (p *Pawn) GetSpawnTime() int64 {
	p.Lock()
	defer p.Unlock()

	return p.spawnTime
}

func (f *Field) SpawnBooster() {
	if len(f.GetBoosters()) >= limitBoosters {
		return
	}

	var spawnPointX int32 = 0
	var spawnPointY int32 = 0

	for {
	repeat:
		spawnPointX = RandomInt32MinMaxN(xBorderLower, xBorderHigher)
		spawnPointY = RandomInt32MinMaxN(yBorderLower, yBorderHigher)

		for _, v := range f.snake.GetCells() {
			if v.positionX == spawnPointX && v.positionY == spawnPointY {
				goto repeat
			}
		}

		for _, v := range f.GetEnemies() {
			if v.positionX == spawnPointX && v.positionY == spawnPointY {
				goto repeat
			}
		}

		break
	}

	b := new(Booster)

	b.positionX = spawnPointX
	b.positionY = spawnPointY
	b.isActive = true
	b.spawnTime = time.Now().Unix()

	f.boosters = append(f.boosters, b)

	f.logger.WriteString(fmt.Sprintf("Spawned booster, %v", b))

	f.Lock()
	f.buffer[b.positionX][b.positionY] = BoosterSymbol
	f.Unlock()
}

func (f *Field) GetBoosters() []*Booster {
	f.Lock()
	defer f.Unlock()

	return f.boosters
}

func (f *Field) DeSpawnBooster(i int) {
	f.Lock()
	defer f.Unlock()

	f.boosters = append(f.boosters[:i], f.boosters[i+1:]...)
}

func (f *Field) DeSpawnEnemy(i int) {
	f.Lock()
	defer f.Unlock()

	e := f.enemies[i]
	f.buffer[e.positionX][e.positionY] = fieldSymbol

	f.enemies = append(f.enemies[:i], f.enemies[i+1:]...)
}

func (f *Field) GetEnemies() []*Enemy {
	f.Lock()
	defer f.Unlock()

	return f.enemies
}

func (f *Field) SpawnEnemy() {
	f.Lock()
	defer f.Unlock()

	if len(f.enemies) >= limitEnemies {
		return
	}

	var spawnPointX int32 = 0
	var spawnPointY int32 = 0

	for {
	repeat:
		spawnPointX = RandomInt32MinMaxN(xBorderLower, xBorderHigher)
		spawnPointY = RandomInt32MinMaxN(yBorderLower, yBorderHigher)

		if f.snake.head.positionX == spawnPointX && f.snake.head.positionY == spawnPointY {
			goto repeat
		}

		for _, v := range f.snake.cells {
			if v.positionX == spawnPointX && v.positionY == spawnPointY {
				goto repeat
			}
		}

		for _, v := range f.boosters {
			if v.positionX == spawnPointX && v.positionY == spawnPointY {
				goto repeat
			}
		}

		break
	}

	e := new(Enemy)

	e.positionX = spawnPointX
	e.positionY = spawnPointY
	e.isActive = true
	e.symbol = EnemyOneShotSymbol
	e.spawnTime = time.Now().Unix()
	e.isOneShot = ItoB(RandomInt32MinMaxI(0, 1))

	if e.isOneShot == false {
		e.countOfCells = RandomInt32MinMaxI(1, 3)
		e.symbol = EnemySymbol
	}

	f.enemies = append(f.enemies, e)

	f.logger.WriteString(fmt.Sprintf("Spawned enemy, %v", e))

	f.buffer[e.positionX][e.positionY] = e.symbol
}
