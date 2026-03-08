package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func FlushAll(req *request.Request) {
	kv.FlushAll()

	req.Client.Write(resp.String("OK"))
}
