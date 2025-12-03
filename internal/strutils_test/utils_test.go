package strutils_test

import (
	"testing"

	"contactsAI/contacts/internal/strutils"

	"github.com/stretchr/testify/assert"
)

func TestPointStr(t *testing.T) {
	inputStr := "test"
	result := strutils.PointStr(inputStr)
	expectedResult := &inputStr
	assert.Exactlyf(t, &inputStr, result, "Result should be equal: %s instead of: %s",
		*expectedResult, *result)
}
