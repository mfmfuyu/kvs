package cmd

import (
	"example.com/kvs/resp"
)

func Ping(args []resp.Value) resp.Value {
	if len(args) > 1 {
		return InvalidArg("ping")
	}

	if len(args) == 0 {
		return resp.String("PONG")
	}

	return resp.String(args[0].Bulk)
}
