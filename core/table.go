package core

import (
	"crypto/md5"
)

type Table struct {
	Ns   string
	Data map[HashKey]uint32 `json:"data"`
	Hash HashKeyMap         `json:"hash"`
}

func NewTable(ns string) *Table {
	return &Table{Ns: ns, Data: make(map[HashKey]uint32), Hash: make(HashKeyMap)}
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

func (t *Table) Get(key []byte) uint32 {
	k := md5.Sum(key)
	_, exists := t.Hash[k]
	if exists {
		return t.Data[k]
	}
	return 0
}

func (t *Table) Destroy() {
	t.Data = nil
	t.Hash = nil
}
