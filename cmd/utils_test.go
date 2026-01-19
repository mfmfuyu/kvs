package cmd

import (
	"testing"

	"example.com/kvs/resp"
	"github.com/stretchr/testify/assert"
)

func TestHasExpiry(t *testing.T) {
	expected := resp.Value{
		Typ: "error",
		Str: "ERR wrong number of arguments for 'test' command",
	}

	assert.Equal(t, expected, InvalidArg("test"))
}
