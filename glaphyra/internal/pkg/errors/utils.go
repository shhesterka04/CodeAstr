package errors

import "github.com/pkg/errors"

func As[T any](err error) bool {
	var target *T
	return errors.As(err, &target)
}
