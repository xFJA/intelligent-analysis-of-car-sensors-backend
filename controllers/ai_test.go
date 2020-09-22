package controllers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	stringToSplit := `["1","2","3","4","5"]`
	expected := []string{"1", "2", "3", "4", "5"}

	result := strings.FieldsFunc(stringToSplit, split)

	assert.Equal(t, expected, result)
}
