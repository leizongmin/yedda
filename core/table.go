package core

import "crypto/md5"

type HashKey [16]byte

type Table struct {
	Data map[HashKey]uint32
	Hash map[HashKey][]byte
}

func NewTable() *Table {
	return &Table{Data: make(map[HashKey]uint32), Hash: make(map[HashKey][]byte)}
}

func (t *Table) Incr(key []byte, n uint32) uint32 {
	k := md5.Sum(key)
	_, exists := t.Hash[k]
	if exists {
		t.Data[k] += n
		return t.Data[k]
	} else {
		t.Hash[k] = key
		t.Data[k] = n
		return n
	}
}

func (t *Table) Destroy() {
	t.Data = nil
	t.Hash = nil
}
