package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
)

func Ttl(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return InvalidArg("ttl")
	}

	key := args[0].Bulk
	ttl := kv.Ttl(key)

	return resp.Integer(ttl)
}
