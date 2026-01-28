package pubsub

import (
	"maps"
	"sync"

	"example.com/kvs/resp"
	"example.com/kvs/server/client"
)

var channelSubscribers = make(map[string]map[*client.Client]struct{})
var clientSubscriptions = make(map[*client.Client]map[string]struct{})

var mutex sync.RWMutex

func Subscribe(c *client.Client, channel string) int {
	mutex.Lock()
	defer mutex.Unlock()

	subscribers, ok := channelSubscribers[channel]
	if !ok {
		// no subscribers
		subscribers = map[*client.Client]struct{}{}
	}

	subscriptions, ok := clientSubscriptions[c]
	if !ok {
		// no subscriptions
		subscriptions = map[string]struct{}{}
	}

	_, ok = subscribers[c]
	if ok {
		// already subscribed
		return len(subscriptions)
	}

	// subscribers
	subscribers[c] = struct{}{}
	channelSubscribers[channel] = subscribers

	// subscriptions
	subscriptions[channel] = struct{}{}
	clientSubscriptions[c] = subscriptions

	return len(subscriptions)
}

func Unsubscribe(c *client.Client, channel string) int {
	mutex.Lock()
	defer mutex.Unlock()

	subscriptions, ok := clientSubscriptions[c]
	if !ok {
		// no subscriptions
		return 0
	}

	subscribers, ok := channelSubscribers[channel]
	if !ok {
		// not subscribed
		return len(subscriptions)
	}

	// subscribers
	delete(subscribers, c)
	channelSubscribers[channel] = subscribers

	// subscriptions
	delete(subscriptions, channel)
	clientSubscriptions[c] = subscriptions

	return len(subscriptions)
}

func Publish(channel string, message string) int {
	mutex.RLock()

	a, ok := channelSubscribers[channel]
	if !ok {
		mutex.RUnlock()
		return 0
	}

	clients := maps.Keys(a)

	mutex.RUnlock()

	msg := resp.Array([]resp.Value{
		resp.Bulk("message"),
		resp.Bulk(channel),
		resp.Bulk(message),
	})

	sent := 0
	for client := range clients {
		client.Write(msg)
		sent++
	}

	return sent
}
