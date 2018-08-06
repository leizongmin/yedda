package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s, err := NewServer(Options{})
	assert.Equal(t, nil, err)
	go s.Loop()
	time.Sleep(2 * time.Second)
	s.Close()
}
