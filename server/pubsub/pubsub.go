package pubsub

import (
	"sync"

	"example.com/kvs/resp"
	"example.com/kvs/server/client"
	"github.com/gobwas/glob"
)

type PubSub struct {
	channelSubscribers  map[string]map[*client.Client]struct{}
	clientSubscriptions map[*client.Client]map[string]struct{}
	patternSubscribers  map[string]map[*client.Client]struct{}
	clientPatterns      map[*client.Client]map[string]struct{}

	mutex sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		channelSubscribers:  make(map[string]map[*client.Client]struct{}),
		clientSubscriptions: make(map[*client.Client]map[string]struct{}),
		patternSubscribers:  make(map[string]map[*client.Client]struct{}),
		clientPatterns:      map[*client.Client]map[string]struct{}{},
	}
}

func (ps *PubSub) internalSubscribe(
	clientToTargets map[*client.Client]map[string]struct{},
	targetToClients map[string]map[*client.Client]struct{},
	c *client.Client,
	t string,
) int {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	subscribers, ok := targetToClients[t]
	if !ok {
		subscribers = map[*client.Client]struct{}{}
	}

	subscriptions, ok := clientToTargets[c]
	if !ok {
		subscriptions = map[string]struct{}{}
	}

	subscribers[c] = struct{}{}
	targetToClients[t] = subscribers

	subscriptions[t] = struct{}{}
	clientToTargets[c] = subscriptions

	return len(ps.clientPatterns[c]) + len(ps.clientSubscriptions[c])
}

func (ps *PubSub) internalUnsubscribe(
	clientToTargets map[*client.Client]map[string]struct{},
	targetToClients map[string]map[*client.Client]struct{},
	c *client.Client,
	t string,
) int {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	subscriptions := clientToTargets[c]
	subscribers := targetToClients[t]

	delete(subscribers, c)
	if len(subscribers) == 0 {
		delete(targetToClients, t)
	}

	delete(subscriptions, t)
	if len(subscriptions) == 0 {
		delete(clientToTargets, c)
	}

	return len(ps.clientPatterns[c]) + len(ps.clientSubscriptions[c])
}

func (ps *PubSub) Subscribe(c *client.Client, channel string) int {
	return ps.internalSubscribe(ps.clientSubscriptions, ps.channelSubscribers, c, channel)
}

func (ps *PubSub) Unsubscribe(c *client.Client, channel string) int {
	return ps.internalUnsubscribe(ps.clientSubscriptions, ps.channelSubscribers, c, channel)
}

func (ps *PubSub) PSubscribe(c *client.Client, channel string) int {
	return ps.internalSubscribe(ps.clientPatterns, ps.patternSubscribers, c, channel)
}

func (ps *PubSub) PUnsubscribe(c *client.Client, channel string) int {
	return ps.internalUnsubscribe(ps.clientPatterns, ps.patternSubscribers, c, channel)
}

func (ps *PubSub) Publish(channel string, message string) int {
	ps.mutex.RLock()

	type patternDelivery struct {
		c *client.Client
		p string
	}
	patternDeliveries := []patternDelivery{}

	for pattern, subscribers := range ps.patternSubscribers {
		g := glob.MustCompile(pattern)
		if g.Match(channel) {
			for subscriber := range subscribers {
				patternDeliveries = append(patternDeliveries, patternDelivery{
					c: subscriber,
					p: pattern,
				})
			}
		}
	}

	deliveries := []*client.Client{}
	subscribers := ps.channelSubscribers[channel]

	for subscriber := range subscribers {
		deliveries = append(deliveries, subscriber)
	}

	ps.mutex.RUnlock()

	msg := resp.Array([]resp.Value{
		resp.Bulk("message"),
		resp.Bulk(channel),
		resp.Bulk(message),
	})

	sent := 0
	for _, delivery := range deliveries {
		delivery.Write(msg)
		sent++
	}

	for _, delivery := range patternDeliveries {
		msg = resp.Array([]resp.Value{
			resp.Bulk("pmessage"),
			resp.Bulk(delivery.p),
			resp.Bulk(channel),
			resp.Bulk(message),
		})
		delivery.c.Write(msg)
		sent++
	}

	return sent
}
