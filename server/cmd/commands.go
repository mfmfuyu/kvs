package cmd

import p "example.com/kvs/server/pubsub"

type Commands struct {
	pubsub *p.PubSub
}

func NewCommands(ps *p.PubSub) *Commands {
	return &Commands{
		pubsub: ps,
	}
}
