package cmd

import (
	"fmt"

	"example.com/kvs/resp"
)

func InvalidArg(cmd string) resp.Value {
	msg := fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd)
	return resp.Error(msg)
}
