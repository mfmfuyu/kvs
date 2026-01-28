package cmd

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/pubsub"
	"example.com/kvs/server/request"
)

func UnSubscribe(req *request.Request) {
	if len(req.Args) < 1 {
		req.Client.Write(InvalidArg("unsubscribe"))
		return
	}

	for _, a := range req.Args {
		channel := a.Bulk
		subscribes := pubsub.Unsubscribe(req.Client, channel)

		req.Client.Write(resp.Array([]resp.Value{
			resp.Bulk("unsubscribe"),
			resp.Bulk(channel),
			resp.Integer(int64(subscribes)),
		}))
	}
}
