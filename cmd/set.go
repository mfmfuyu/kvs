package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
)

func Set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return InvalidArg("set")
	}

	key := args[0].Bulk
	value := args[1].Bulk

	kv.Set(key, value)

	return resp.String(value)
}
