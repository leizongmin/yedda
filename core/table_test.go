package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable("")
	assert.NotEqual(t, nil, table)
}

func TestTable_Incr(t *testing.T) {
	table := NewTable("")
	assert.Equal(t, uint32(1), table.Incr([]byte("abc"), 1))
	assert.Equal(t, uint32(2), table.Incr([]byte("abc"), 1))
	assert.Equal(t, uint32(3), table.Incr([]byte("abc"), 1))
	assert.Equal(t, uint32(4), table.Incr([]byte("abc"), 1))
	assert.Equal(t, uint32(5), table.Incr([]byte("abc"), 1))
}

func TestTable_Destroy(t *testing.T) {
	table := NewTable("")
	table.Incr([]byte("666"), 1)
	table.Incr([]byte("abcdefg"), 1)
	table.Destroy()
	assert.Equal(t, 0, len(table.Hash))
	assert.Equal(t, 0, len(table.Data))
}
