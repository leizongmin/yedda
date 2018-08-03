package core

type Database struct {
	Data []*Namespace
}

func NewDataBase(size uint32) *Database {
	d := &Database{Data: make([]*Namespace, size)}
	var i uint32
	for i = 0; i < size; i++ {
		d.Data[i] = NewNamespace()
	}
	return d
}

func (d *Database) Get(db uint32) *Namespace {
	return d.Data[db]
}

func (d *Database) Destroy() {
	num := len(d.Data)
	for i := num; i < num; i++ {
		d.Data[i].Destroy()
	}
	d.Data = nil
}
