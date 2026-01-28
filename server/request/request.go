package request

import (
	"example.com/kvs/resp"
	"example.com/kvs/server/client"
)

type Request struct {
	Args   []resp.Value
	Client *client.Client
}
