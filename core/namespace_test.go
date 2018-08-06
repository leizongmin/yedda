package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNamespace(t *testing.T) {
	n := NewNamespace(0)
	assert.NotEqual(t, nil, n)
}

func TestNamespace_GetTable(t *testing.T) {
	n := NewNamespace(0)
	a := n.Get("abc", 1)
	assert.NotEqual(t, nil, a)

	b := n.Get("efg", 1)
	assert.NotEqual(t, nil, b)
	assert.NotEqual(t, b, a)

	assert.Equal(t, a, n.Get("abc", 1))

	a.Incr([]byte("11111"), 1)
	b.Incr([]byte("11111"), 1)
}

func TestNamespace_Destroy(t *testing.T) {
	n := NewNamespace(0)
	n.Get("a", 1)
	n.Get("b", 1)
	n.Get("c", 1)
	n.Destroy()
	assert.Equal(t, 0, len(n.Data))
}
