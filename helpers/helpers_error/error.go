package helpers_error

import (
	"github.com/pkg/errors"
	"strings"
)

func PanicIfError(e error, info ...string) {
	if e != nil {
		panic(WrapError(e, info...))
	}
}

func WrapError(e error, info ...string) error {
	return errors.Wrap(e, strings.Join(info, " --> "))
}
