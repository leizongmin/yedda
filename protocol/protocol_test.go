package protocol

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPackage_Pack(t *testing.T) {
	p := NewPackage(1, OpPing, []byte("abc"))
	buf := bytes.NewBuffer(make([]byte, 0))
	err := p.Pack(buf)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{0, 1, 0, 1, 0, 3, 97, 98, 99}, buf.Bytes())
}

func TestNewPackageFromReader(t *testing.T) {
	p := NewPackage(2, OpPong, []byte("123456"))
	buf := bytes.NewBuffer(make([]byte, 0))
	err := p.Pack(buf)
	assert.Equal(t, nil, err)
	p2, err := NewPackageFromReader(buf)
	assert.Equal(t, nil, err)
	assert.Equal(t, p, p2)
}

func BenchmarkPackage_Pack(b *testing.B) {
	p := NewPackage(1, OpPing, []byte("abc"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Pack(bytes.NewBuffer(make([]byte, 0)))
	}
}

func BenchmarkNewPackageFromReader(b *testing.B) {
	p := NewPackage(1, OpPing, []byte("abc"))
	buf := bytes.NewBuffer(make([]byte, 0))
	p.Pack(buf)
	buf2 := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewPackageFromReader(bytes.NewBuffer(buf2))
	}
}
