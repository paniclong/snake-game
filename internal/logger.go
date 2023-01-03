package internal

import (
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func (l *Logger) Init() error {
	file, err := os.OpenFile("./test.log", os.O_APPEND|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return err
	}

	l.file = file

	err = l.WriteString("Init logger")

	if err != nil {
		return err
	}

	return nil
}

func (l *Logger) Close() error {
	err := l.file.Close()

	if err != nil {
		return err
	}

	return nil
}

func (l *Logger) WriteString(s string, disabled ...bool) error {
	if len(disabled) > 0 && disabled[0] {
		return nil
	}

	t := time.Now().Format(time.RFC3339Nano)

	_, err := l.file.WriteString("[" + t + "] " + s + "\n")
	if err != nil {
		return err
	}

	return nil
}
