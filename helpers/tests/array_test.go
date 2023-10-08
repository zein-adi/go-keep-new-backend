package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"testing"
)

func TestSlice(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, helpers.Slice(data, 0, 5))
	assert.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9}, helpers.Slice(data, 1, 8))
	assert.Equal(t, []int{6, 7, 8, 9}, helpers.Slice(data, 5, 5))
	assert.Equal(t, []int{}, helpers.Slice(data, 10, 5))
}
