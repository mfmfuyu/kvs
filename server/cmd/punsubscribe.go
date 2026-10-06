package cmd

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func (c *Commands) PUnsubscribe(req *request.Request) {
	if len(req.Args) < 1 {
		req.Client.Write(InvalidArg("punsubscribe"))
		return
	}

	for _, a := range req.Args {
		pattern := a.Bulk
		subscribes := c.pubsub.PUnsubscribe(req.Client, pattern)

		req.Client.Write(resp.Array([]resp.Value{
			resp.Bulk("punsubscribe"),
			resp.Bulk(pattern),
			resp.Integer(int64(subscribes)),
		}))
	}
}
