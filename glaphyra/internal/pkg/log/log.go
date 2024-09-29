package log

import (
	"log"

	"github.com/pkg/errors"
)

func WrapErr(err error) error {
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, getCaller())
	log.Println(err.Error())
	return err
}

func Wrap(err error) error {
	return errors.Wrap(err, getCaller())
}
