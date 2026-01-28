package cmd

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func Ping(req *request.Request) {
	if len(req.Args) > 1 {
		req.Client.Write(InvalidArg("ping"))
		return
	}

	if len(req.Args) == 0 {
		req.Client.Write(resp.String("PONG"))
		return
	}

	req.Client.Write(resp.String(req.Args[0].Bulk))
}
