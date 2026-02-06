package kv

import (
	"math"
	"sync"
	"time"
)

type Object struct {
	Value     string
	ExpiresAt time.Time
}

func (o Object) HasExpiry() bool {
	return !o.ExpiresAt.IsZero()
}

func (o Object) IsExpired(t time.Time) bool {
	if !o.HasExpiry() {
		return false
	}

	return !t.Before(o.ExpiresAt)
}

var objects = map[string]Object{}
var expires = map[string]struct{}{}

var mutex = sync.RWMutex{}

func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	obj, ok := objects[key]
	if !ok {
		return "", false
	}

	if obj.IsExpired(time.Now()) {
		return "", false
	}

	return obj.Value, true
}

func Set(key string, value string) {
	mutex.Lock()
	defer mutex.Unlock()

	objects[key] = Object{
		Value: value,
	}
}

func SetExpiresAt(key string, expireAt time.Time) bool {
	mutex.Lock()
	defer mutex.Unlock()

	object, ok := objects[key]
	if !ok {
		return false
	}

	object.ExpiresAt = expireAt
	objects[key] = object

	expires[key] = struct{}{}

	return true
}

func Del(key string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(objects, key)
	delete(expires, key)
}

func Ttl(key string) int64 {
	mutex.RLock()
	defer mutex.RUnlock()

	obj, ok := objects[key]
	if !ok {
		return -2
	}

	if !obj.HasExpiry() {
		return -1
	}

	diff := time.Until(obj.ExpiresAt)
	if diff < 0 {
		return 0
	}

	return int64(math.Ceil(diff.Seconds()))
}

func ActiveExpire() {
	mutex.Lock()
	defer mutex.Unlock()

	budget := 10

	for key := range expires {
		budget--
		if budget < 0 {
			break
		}

		obj, ok := objects[key]
		if !ok {
			continue
		}

		if obj.IsExpired(time.Now()) {
			delete(objects, key)
			delete(expires, key)
		}
	}
}
