package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func Get(req *request.Request) {
	if len(req.Args) != 1 {
		req.Client.Write(InvalidArg("get"))
		return
	}

	key := req.Args[0].Bulk
	value, ok := kv.Get(key)
	if !ok {
		req.Client.Write(resp.Null())
		return
	}

	req.Client.Write(resp.Bulk(value))
	return
}
