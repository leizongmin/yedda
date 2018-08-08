package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCmdArg(t *testing.T) {
	a := NewCmdArg(1, "abc", 123, []byte("efg456"), 99)
	b, err := a.Bytes()
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{0, 0, 0, 1, 3, 97, 98, 99, 0, 0, 0, 123, 6, 101, 102, 103, 52, 53, 54, 0, 0, 0, 99}, b)
	a2, err := NewCmdArgFromBytes(b)
	assert.Equal(t, nil, err)
	assert.Equal(t, a, a2)
}

func BenchmarkNewCmdArg_Bytes(b *testing.B) {
	a := NewCmdArg(1, "abc", 123, []byte("efg456"), 99)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Bytes()
	}
}

func BenchmarkNewCmdArgFromBytes(b *testing.B) {
	a := NewCmdArg(1, "abc", 123, []byte("efg456"), 99)
	buf, _ := a.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewCmdArgFromBytes(buf)
	}
}
