package server

import (
	"github.com/leizongmin/simple-limiter-service/client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer(Options{DatabaseSize: 2, EnableLog: true})
	assert.Equal(t, nil, err)
	go s.Loop()

	testClient := func(db uint32) {
		c, err := client.NewClient(client.Options{Db: db})
		defer c.Close()
		assert.Equal(t, nil, err)

		v1, err := c.Incr("hello", "world", 100)
		assert.Equal(t, nil, err)
		assert.Equal(t, uint32(1), v1)

		v2, err := c.Incr("hello", "world", 100)
		assert.Equal(t, nil, err)
		assert.Equal(t, uint32(2), v2)

		v3, err := c.Get("hello", "world2", 100)
		assert.Equal(t, nil, err)
		assert.Equal(t, uint32(0), v3)

		v4, err := c.Get("hello", "world", 100)
		assert.Equal(t, nil, err)
		assert.Equal(t, uint32(2), v4)

		v5, err := c.Ping()
		assert.Equal(t, nil, err)
		assert.Equal(t, true, v5 >= 0)
	}
	go testClient(0)
	go testClient(1)

	time.Sleep(2 * time.Second)
	s.Close()
}

func TestServer_Many(t *testing.T) {
	s, err := NewServer(Options{DatabaseSize: 2, EnableLog: true})
	defer s.Close()
	assert.Equal(t, nil, err)
	go s.Loop()

	c, err := client.NewClient(client.Options{Db: 1})
	defer c.Close()
	assert.Equal(t, nil, err)

	for i := 0; i < 100; i++ {
		v, err := c.Incr("aaa", "abc", 100)
		assert.Equal(t, nil, err)
		v2, err := c.Get("aaa", "abc", 100)
		assert.Equal(t, nil, err)
		assert.Equal(t, v, v2)
		assert.Equal(t, true, v > 0)
	}
}

func BenchmarkNewServer_Incr(b *testing.B) {
	s, _ := NewServer(Options{})
	defer s.Close()
	go s.Loop()
	c, _ := client.NewClient(client.Options{})
	defer c.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Incr("abc", "1234", 100)
	}
}

func BenchmarkNewServer_Get(b *testing.B) {
	s, _ := NewServer(Options{})
	defer s.Close()
	go s.Loop()
	c, _ := client.NewClient(client.Options{})
	defer c.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("abc", "1234", 100)
	}
}
