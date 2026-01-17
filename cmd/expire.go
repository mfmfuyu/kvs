package cmd

import (
	"strconv"
	"time"

	"example.com/kvs/kv"
	"example.com/kvs/resp"
)

func Expire(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return InvalidArg("expire")
	}

	key := args[0].Bulk
	rawSeconds := args[1].Bulk
	seconds, err := strconv.ParseInt(rawSeconds, 10, 64)
	if err != nil {
		return resp.Error(ERR_NUM)
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(seconds) * time.Second)

	ok := kv.SetExpiresAt(key, expiresAt)
	if !ok {
		return resp.Integer(0)
	}

	return resp.Integer(1)
}
