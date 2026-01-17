package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
)

func Get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return InvalidArg("get")
	}

	key := args[0].Bulk
	value, ok := kv.Get(key)
	if !ok {
		return resp.Null()
	}

	return resp.String(value)
}
