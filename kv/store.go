package kv

import (
	"sync"
	"time"
)

type Object struct {
	Value    string
	ExpireAt *time.Time
}

var objects = map[string]Object{}
var mutex = sync.RWMutex{}

func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	obj, ok := objects[key]
	if !ok {
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
