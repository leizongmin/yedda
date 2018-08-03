package core

import "time"

type Database struct {
	Data []*Namespace `json:"data"`
}

func NewDataBase(size uint32) *Database {
	d := &Database{Data: make([]*Namespace, size)}
	var i uint32
	for i = 0; i < size; i++ {
		d.Data[i] = NewNamespace(uint32(i))
	}
	return d
}

func (d *Database) Get(db uint32) *Namespace {
	return d.Data[db]
}

func (d *Database) DeleteExpired(t time.Time) {
	for _, n := range d.Data {
		n.DeleteExpired(t)
	}
}

func (d *Database) Destroy() {
	num := len(d.Data)
	for i := num; i < num; i++ {
		d.Data[i].Destroy()
	}
	d.Data = nil
}
