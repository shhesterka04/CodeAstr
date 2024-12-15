package log

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func WriteLog(msg string) {
	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}

func Error(err error) {
	if err == nil {
		return
	}
	err = errors.Wrap(err, getCaller())
	WriteLog(err.Error())
}

func WrapErr(err error) error {
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, getCaller())
	WriteLog(err.Error())
	return err
}

func Wrap(err error) error {
	return errors.Wrap(err, getCaller())
}
