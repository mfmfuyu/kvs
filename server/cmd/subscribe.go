package cmd

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/pubsub"
	"example.com/kvs/server/request"
)

func Subscribe(req *request.Request) {
	if len(req.Args) < 1 {
		req.Client.Write(InvalidArg("subscribe"))
		return
	}

	for _, a := range req.Args {
		channel := a.Bulk
		subscribes := pubsub.Subscribe(req.Client, channel)

		req.Client.Write(resp.Array([]resp.Value{
			resp.Bulk("subscribe"),
			resp.Bulk(channel),
			resp.Integer(int64(subscribes)),
		}))
	}
}
