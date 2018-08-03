package core

import "testing"

func TestNewNamespace(t *testing.T) {
	n := NewNamespace()
	if n == nil {
		t.Error("unexpected nil Namespace pointer")
	}
}

func TestNamespace_GetTable(t *testing.T) {
	n := NewNamespace()
	a := n.Get("abc")
	if a == nil {
		t.Error("unexpected nil Table pointer")
	}
	b := n.Get("efg")
	if b == nil {
		t.Error("unexpected nil Table pointer")
	}
	if a == b {
		t.Error("get different table name")
	}
	if a != n.Get("abc") {
		t.Error("get the same table name")
	}
	a.Incr([]byte("11111"), 1)
	b.Incr([]byte("11111"), 1)
}

func TestNamespace_Destroy(t *testing.T) {
	n := NewNamespace()
	n.Get("a")
	n.Get("b")
	n.Get("c")
	n.Destroy()
	if n.Data != nil {
		t.Errorf("expected namespace.Data nil pointer")
	}
}
