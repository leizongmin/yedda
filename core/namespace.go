package core

type Namespace struct {
	Data map[string]*Table
}

func NewNamespace() *Namespace {
	return &Namespace{Data: make(map[string]*Table)}
}

func (n *Namespace) Get(name string) *Table {
	_, exists := n.Data[name]
	if !exists {
		n.Data[name] = NewTable()
	}
	return n.Data[name]
}

func (n *Namespace) Destroy() {
	for k := range n.Data {
		n.Data[k].Destroy()
	}
	n.Data = nil
}
