package kv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHasExpiry(t *testing.T) {
	object := Object{}
	assert.False(t, object.HasExpiry())

	object = Object{
		ExpiresAt: time.Now(),
	}
	assert.True(t, object.HasExpiry())
}

func TestIsExpired(t *testing.T) {
	now := time.Now()

	object := Object{}
	assert.False(t, object.IsExpired(now))

	expiresAt := now.Add(1 * time.Minute)

	object = Object{
		ExpiresAt: expiresAt,
	}

	assert.False(t, object.IsExpired(now))
	assert.True(t, object.IsExpired(expiresAt))

	now2 := now.Add(2 * time.Minute)
	assert.True(t, object.IsExpired(now2))
}
