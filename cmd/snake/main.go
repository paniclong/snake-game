package main

import (
	"github.com/paniclong/snake-game/internal"
	"log"
	"math/rand"
	"time"
)

/*
#include <curses.h>
#include <stdio.h>
#cgo LDFLAGS: -lcurses
*/
import "C"

const baseDelay = 700
const minDelay = 100

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := *new(internal.Logger)
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}

	snake := *new(internal.Snake)
	snake.Init()

	field := *new(internal.Field)
	field.Init(&snake, &logger)

	defer logger.Close()

	// Init screen
	C.initscr()
	C.cbreak()
	C.noecho()
	// Allows using F1...F12, arrows and other keys
	C.keypad(C.stdscr, true)

	defer C.endwin()

	go field.ChangeDirectionByKey()

	// Base delay is 300 mc
	// Each new cell on the snake minus 5 mc, but not less than 100 mc
	delay := baseDelay * time.Millisecond
	sleepTime := delay

	for {
		C.clear()

		if !field.IsActive() {
			C.refresh()

			return
		}

		field.SpawnBooster()
		err := field.OnStep()
		if err != nil {
			return
		}

		field.Print()
		C.refresh()

		size := field.GetSnake().GetSize()

		if sleepTime < minDelay {
			sleepTime = minDelay
		} else {
			sleepTime = delay - (time.Duration(size*5) * time.Millisecond)
		}

		time.Sleep(sleepTime)
	}
}
