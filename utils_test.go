package getaredis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	tests := [...]int{0, 10, 20, 30}

	for test := range tests {
		str := generateRandomString(test)
		assert.Equal(t, test, len(str), "String should be equal the specified length")
	}
}
