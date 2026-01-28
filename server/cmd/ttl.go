package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func Ttl(req *request.Request) {
	if len(req.Args) != 1 {
		req.Client.Write(InvalidArg("ttl"))
		return
	}

	key := req.Args[0].Bulk
	ttl := kv.Ttl(key)

	req.Client.Write(resp.Integer(ttl))
}
