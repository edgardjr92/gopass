package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		// given
		arr := []int{1, 2, 3}

		// when
		result := Map(arr, func(n int) int {
			return n * 2
		})

		// then
		expected := []int{2, 4, 6}
		assert.Equal(t, expected, result)
	})

	t.Run("string", func(t *testing.T) {
		// given
		arr := []string{"a", "b", "c"}

		// when
		result := Map(arr, func(n string) string {
			return n + "!"
		})

		// then
		expected := []string{"a!", "b!", "c!"}
		assert.Equal(t, expected, result)
	})
}

func TestIsBlank(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		expected bool
	}{
		{"empty string", "", true},
		{"blank string", " ", true},
		{"not blank string", "a", false},
		{"not blank string with space", " a", false},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			// when
			result := IsBlank(tc.str)

			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}
