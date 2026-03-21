package imagegen

import "errors"

func errNotImplemented(msg string) error {
	return errors.New(msg)
}
