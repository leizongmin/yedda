package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDataBase(t *testing.T) {
	d := NewDataBase(16)
	assert.NotEqual(t, nil, d)
	assert.Equal(t, 16, len(d.Data))

	d2 := NewDataBase(32)
	assert.NotEqual(t, nil, d2)
	assert.Equal(t, 32, len(d2.Data))

	d3 := NewDataBase(512)
	assert.NotEqual(t, nil, d3)
	assert.Equal(t, 512, len(d3.Data))
}
