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

	err = l.writeString("Init logger")

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

func (l *Logger) writeString(s string) error {
	t := time.Now().Format(time.RFC3339Nano)

	_, err := l.file.WriteString("[" + t + "] " + s + "\n")
	if err != nil {
		return err
	}

	return nil
}
