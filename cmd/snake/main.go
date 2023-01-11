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

const baseDelay = 300
const minDelay = 100

const exitFailureCode = 0x1

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := *new(internal.Logger)
	err := logger.Init()
	if err != nil {
		log.Panic(err)
	}

	defer func(l *internal.Logger) {
		defer l.Close()

		if r := recover(); r != nil {
			l.WriteString(fmt.Sprint("Something went wrong: ", r))

			os.Exit(exitFailureCode)
		}
	}(&logger)

	snake := *new(internal.Snake)
	snake.Init()

	field := *new(internal.Field)
	field.Init(&snake, &logger)

	CatchSignals(&field)

	// Init screen
	C.initscr()
	C.cbreak()
	C.noecho()
	// Allows using F1...F12, arrows and other keys
	C.keypad(C.stdscr, true)

	defer C.endwin()

	C.addstr(C.CString(GetMenuText()))
	C.addstr(C.CString("Press any key for start game"))
	C.getch()

	go field.ChangeDirectionByKey()

	// Base delay is 300 mc
	// Each new cell on the snake minus 5 mc, but not less than 100 mc
	delay := baseDelay * time.Millisecond
	sleepTime := delay

	for {
		if !field.IsActive() {
			C.refresh()

			return
		}

		if !field.IsFirstStart {
			err = field.OnStep()
			if err != nil {
				logger.WriteString(fmt.Sprint("On step error: ", err.Error()))
				C.addstr(C.CString(err.Error()))

				return
			}
		} else {
			field.IsFirstStart = false
		}

		field.ReCalcBoosters()
		field.ReCalcEnemies()

		C.clear()
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
	return fmt.Sprintf("Symbols: \n 1. %s - snake head \n 2. %s - snake body \n 3. %s - boosters \n "+
		"4. %s - enemy, one shot \n 5. %s - enemy, no one shot \n",
		string(internal.HeadSymbol),
		string(internal.BodySymbol),
		string(internal.BoosterSymbol),
		string(internal.EnemyOneShotSymbol),
		string(internal.EnemySymbol),
	)
}
