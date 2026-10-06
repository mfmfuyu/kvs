package cmd

import (
	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func (c *Commands) FlushAll(req *request.Request) {
	kv.FlushAll()

	req.Client.Write(resp.String("OK"))
}
