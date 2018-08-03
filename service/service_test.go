package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	s := NewService(Options{DatabaseSize: 2, TickerDuration: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	for i := 0; i < 5; i++ {
		s.CmdIncr(0, "test1", 100, []byte("a"), 1)
	}
	for i := 0; i < 7; i++ {
		s.CmdIncr(0, "test1", 100, []byte("b"), 1)
	}

	assert.Equal(t, uint32(0), s.CmdGet(0, "test1", 100, []byte("c")))
	assert.Equal(t, uint32(0), s.CmdGet(1, "test1", 100, []byte("c")))

	assert.Equal(t, uint32(0), s.CmdGet(1, "test1", 100, []byte("a")))
	assert.Equal(t, uint32(0), s.CmdGet(1, "test1", 200, []byte("a")))
	assert.Equal(t, uint32(0), s.CmdGet(0, "test1", 200, []byte("a")))
	assert.Equal(t, uint32(5), s.CmdGet(0, "test1", 100, []byte("a")))

	assert.Equal(t, uint32(0), s.CmdGet(1, "test1", 100, []byte("b")))
	assert.Equal(t, uint32(0), s.CmdGet(0, "test2", 100, []byte("b")))
	assert.Equal(t, uint32(7), s.CmdGet(0, "test1", 100, []byte("b")))

	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, uint32(0), s.CmdGet(0, "test1", 100, []byte("a")))
	assert.Equal(t, uint32(0), s.CmdGet(0, "test1", 100, []byte("b")))
	for i := 0; i < 5; i++ {
		s.CmdIncr(0, "test1", 100, []byte("a"), 3)
	}
	assert.Equal(t, uint32(15), s.CmdGet(0, "test1", 100, []byte("a")))

	time.Sleep(110 * time.Millisecond)
	//fmt.Println(s.CmdGet(0, "test1", 100, []byte("a")))
	assert.Equal(t, uint32(0), s.CmdGet(0, "test1", 100, []byte("a")))

	//spew.Dump(s.database.Get(0).Get("test1", 100*time.Millisecond))
	//spew.Dump(s.database.Get(0).Data)
	//n := s.database.Get(0)
	//fmt.Printf("%x end: %+v\n", &n, n.Data)
	//fmt.Println(s.CmdGet(0, "test1", 100, []byte("a")))
}
