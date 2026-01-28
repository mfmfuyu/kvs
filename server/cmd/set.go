package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func Set(req *request.Request) {
	if len(req.Args) != 2 {
		req.Client.Write(InvalidArg("set"))
		return
	}

	key := req.Args[0].Bulk
	value := req.Args[1].Bulk

	kv.Set(key, value)

	req.Client.Write(resp.String("OK"))
}
