package cmd

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func (c *Commands) Publish(req *request.Request) {
	if len(req.Args) != 2 {
		req.Client.Write(InvalidArg("publish"))
		return
	}

	channel := req.Args[0].Bulk
	message := req.Args[1].Bulk

	subscribers := c.pubsub.Publish(channel, message)
	req.Client.Write(resp.Integer(int64(subscribers)))
}
