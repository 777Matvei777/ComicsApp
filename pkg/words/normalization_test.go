package words

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandContractions(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"aren't", "are not"},
		{"can't", "cannot"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ExpandContractions(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSplitString(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"Hello, world!", []string{"Hello", "world", ""}},
		{"It's a test.", []string{"It's", "a", "test", ""}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SplitString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStemming(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{[]string{"running", "jumps", "easily"}, []string{"run", "jump", "easili"}},
		{[]string{"better", "driving", "cars"}, []string{"better", "drive", "car"}},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.input, ","), func(t *testing.T) {
			result, err := Stemming(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
