package main

import (
	"fmt"
	"github.com/paniclong/snake-game/internal"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
#include <curses.h>
#include <stdio.h>
#cgo LDFLAGS: -lcurses
*/
import "C"

// Delay for snake step
const baseDelay = 300
const minDelay = 100

// fps for display game field
const fps = 30

const exitFailureCode = 0x1

func main() {
	rand.Seed(time.Now().UnixNano())

	// Init logger
	logger, err := internal.CreateLogger()
	if err != nil {
		log.Fatalln(err)
	}

	// Init first of defer function for logging all errors (panic, etc.)
	// and for return correct exit code
	defer func(l *internal.Logger) {
		defer l.Close()

		if r := recover(); r != nil {
			l.WriteString(fmt.Sprint("Something went wrong: ", r))

			os.Exit(exitFailureCode)
		}
	}(logger)

	// Init game entities
	snake := internal.CreateSnake()
	field := internal.CreateField(snake, logger)

	// Catching signals for the correct exit from the game
	CatchSignals(field)

	// Init screen
	C.initscr()
	// Read one character, don't wait
	C.cbreak()
	// Don't display input
	C.noecho()
	// Allows using F1...F12, arrows and other keys
	C.keypad(C.stdscr, true)
	// Restore original tty modes
	defer C.endwin()

	C.addstr(C.CString(GetMenuText()))
	C.addstr(C.CString("Press any key for start game"))
	C.getch()

	go field.ChangeDirectionByKey()
	go Step(field, logger)

	// Main game loop
	for {
		if !field.IsActive() {
			C.refresh()

			return
		}

		C.clear()
		field.Print()
		C.touchwin(C.stdscr)
		C.refresh()

		field.ReCalcBoosters()
		field.ReCalcEnemies()

		time.Sleep(time.Second / fps)
	}
}

func CatchSignals(f *internal.Field) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func(f *internal.Field) {
		<-sigChan

		f.SetIsActive(false)
	}(f)
}

func GetMenuText() string {
	return fmt.Sprintf("Symbols: \n 1. %c - snake head \n 2. %c - snake body \n 3. %c - boosters \n "+
		"4. %c - enemy, one shot \n 5. %c - enemy, no one shot \n",
		internal.HeadSymbol,
		internal.BodySymbol,
		internal.BoosterSymbol,
		internal.EnemyOneShotSymbol,
		internal.EnemySymbol,
	)
}

func Step(field *internal.Field, logger *internal.Logger) {
	// Base delay is 300 mc
	// Each new cell on the snake minus 5 mc, but not less than 100 mc
	delay := baseDelay * time.Millisecond
	sleepTime := delay

	for {
		err := field.OnStep()
		if err != nil {
			logger.WriteString(fmt.Sprint("On step error: ", err.Error()))
			C.addstr(C.CString(err.Error()))

			field.SetIsActive(false)

			return
		}

		size := field.GetSnake().GetSize()

		if sleepTime < minDelay {
			sleepTime = minDelay
		} else {
			sleepTime = delay - (time.Duration(size*5) * time.Millisecond)
		}

		time.Sleep(sleepTime)
	}
}
