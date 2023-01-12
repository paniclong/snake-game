package internal

import (
	"fmt"
	"os"
	"time"
)

const logFilePath = "./logs/%s.log"

type Logger struct {
	file *os.File
}

func CreateLogger() (*Logger, error) {
	logger := new(Logger)
	err := logger.Init()
	if err != nil {
		return logger, err
	}

	return logger, nil
}

func (l *Logger) Init() error {
	t := time.Now().Format("2006-01-01")

	file, err := os.OpenFile(fmt.Sprintf(logFilePath, t), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)

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
