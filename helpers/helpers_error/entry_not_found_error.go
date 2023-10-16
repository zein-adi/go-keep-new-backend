package helpers_error

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	EntryNotFoundError      = errors.New("entry not found")
	EntryCountMismatchError = errors.New("entry count mismatch")
)

func NewEntryNotFoundError(entityName string, keyField string, keyValue string) error {
	err := fmt.Sprintf("%s with %s = %s", entityName, keyField, keyValue)
	return errors.Wrap(EntryNotFoundError, err)
}

func NewEntryCountMismatchError(expected int, actual int) error {
	err := fmt.Sprintf("expected %d actual %d", expected, actual)
	return errors.Wrap(EntryCountMismatchError, err)
}
