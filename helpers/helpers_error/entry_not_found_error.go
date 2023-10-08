package helpers_error

import (
	"fmt"
	"github.com/pkg/errors"
)

var EntryNotFoundError = errors.New("entry not found")

func NewEntryNotFoundError(entityName string, keyField string, keyValue string) error {
	err := fmt.Sprintf("%s with %s = %s", entityName, keyField, keyValue)
	return errors.Wrap(EntryNotFoundError, err)
}
