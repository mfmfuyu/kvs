package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type TestCase struct {
		value    string
		expected []string
	}

	var testCases = []TestCase{
		// normal text
		{
			value:    "test",
			expected: []string{"test"},
		},
		// with space
		{
			value:    "te st",
			expected: []string{"te", "st"},
		},
		// with single quote
		{
			value:    "'test'",
			expected: []string{"test"},
		},
		// with double quote
		{
			value:    "\"test\"",
			expected: []string{"test"},
		},
		// with invalid single quote
		{
			value:    "'test",
			expected: nil,
		},
		// with invalid double quote
		{
			value:    "\"test",
			expected: nil,
		},
		// with double quote in single quote
		{
			value:    "'\"test\"'",
			expected: []string{"\"test\""},
		},
		// with single quote in double quote
		{
			value:    "\"'test'\"",
			expected: []string{"'test'"},
		},
		// escape with single quote
		{
			value:    "\\'test",
			expected: []string{"'test"},
		},
		// escape with double quote
		{
			value:    "\\\"test",
			expected: []string{"\"test"},
		},
		// some tests...
		{
			value:    "test 123 'hello world' \"123 456\" \\'test",
			expected: []string{"test", "123", "hello world", "123 456", "'test"},
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, Parse(testCase.value))
	}
}

func TestDigits(t *testing.T) {
	type TestCase struct {
		value    int
		expected int
	}

	var testCases = []TestCase{
		{
			value:    1,
			expected: 1,
		},
		{
			value:    0,
			expected: 1,
		},
		{
			value:    10,
			expected: 2,
		},
		{
			value:    -1,
			expected: 1,
		},
		{
			value:    -10,
			expected: 2,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, Digits(testCase.value))
	}
}
