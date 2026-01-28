package cmd

import (
	"fmt"

	"example.com/kvs/resp"
)

const ERR_NUM = "value is not an integer or out of range"

func InvalidArg(cmd string) resp.Value {
	msg := fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd)
	return resp.Error(msg)
}
