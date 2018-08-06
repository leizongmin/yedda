package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestNewCmdArg(t *testing.T) {
	a := NewCmdArg(1, "abc", 123, []byte("efg456"), 99)
	fmt.Println(a)
	b, err := a.Bytes()
	assert.Equal(t, nil, err)
	fmt.Println(b)
	a2, err := NewCmdArgFromBytes(b)
	assert.Equal(t, nil, err)
	fmt.Println(a2)
}

func TestNewService(t *testing.T) {
	s := NewService(Options{DatabaseSize: 2, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	for i := 0; i < 5; i++ {
		s.Incr(NewCmdArg(0, "test1", 100, []byte("a"), 1))
	}
	for i := 0; i < 7; i++ {
		s.Incr(NewCmdArg(0, "test1", 100, []byte("b"), 1))
	}

	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test1", 100, []byte("c"), 0)))
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(1, "test1", 100, []byte("c"), 0)))

	assert.Equal(t, uint32(0), s.Get(NewCmdArg(1, "test1", 100, []byte("a"), 0)))
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(1, "test1", 200, []byte("a"), 0)))
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test1", 200, []byte("a"), 0)))
	assert.Equal(t, uint32(5), s.Get(NewCmdArg(0, "test1", 100, []byte("a"), 0)))

	assert.Equal(t, uint32(0), s.Get(NewCmdArg(1, "test1", 100, []byte("b"), 0)))
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test2", 100, []byte("b"), 0)))
	assert.Equal(t, uint32(7), s.Get(NewCmdArg(0, "test1", 100, []byte("b"), 0)))

	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test1", 100, []byte("a"), 0)))
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test1", 100, []byte("b"), 0)))
	for i := 0; i < 5; i++ {
		s.Incr(NewCmdArg(0, "test1", 100, []byte("a"), 3))
	}
	assert.Equal(t, uint32(15), s.Get(NewCmdArg(0, "test1", 100, []byte("a"), 0)))

	time.Sleep(110 * time.Millisecond)
	assert.Equal(t, uint32(0), s.Get(NewCmdArg(0, "test1", 100, []byte("a"), 0)))
}

func TestNewServiceParallels(t *testing.T) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	go func() {
		for {
			c := s.Incr(NewCmdArg(0, "a", 50, []byte("abc"), 1))
			assert.Equal(t, true, c > 0)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		}
	}()
	go func() {
		for {
			c := s.Incr(NewCmdArg(0, "b", 50, []byte("abc"), 1))
			assert.Equal(t, true, c > 0)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		}
	}()
	go func() {
		for {
			c := s.Incr(NewCmdArg(0, "a", 60, []byte("abc"), 1))
			assert.Equal(t, true, c > 0)
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
		}
	}()

	time.Sleep(time.Second)
	s.Stop()
}

func BenchmarkService_Incr_1(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "a", 100, []byte("b"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Incr(a)
	}
}

func BenchmarkService_Incr_10(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "abcdefghi", 100, []byte("abcdefghi"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Incr(a)
	}
}

func BenchmarkService_Incr_36(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "abcdefghijklmnopqrstuvwxyz1234567890", 100, []byte("abcdefghijklmnopqrstuvwxyz1234567890"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Incr(a)
	}
}

func BenchmarkService_Get_1(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "a", 100, []byte("b"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(a)
	}
}

func BenchmarkService_Get10(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "abcdefghi", 100, []byte("abcdefghi"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(a)
	}
}

func BenchmarkService_Get_36(b *testing.B) {
	s := NewService(Options{DatabaseSize: 1, TimeAccuracy: 10 * time.Millisecond})
	s.Start()
	defer s.Destroy()

	a := NewCmdArg(0, "abcdefghijklmnopqrstuvwxyz1234567890", 100, []byte("abcdefghijklmnopqrstuvwxyz1234567890"), 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(a)
	}
}
