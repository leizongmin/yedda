package core

import "testing"

func TestNewNamespace(t *testing.T) {
	n := NewNamespace(0)
	if n == nil {
		t.Error("unexpected nil Namespace pointer")
	}
}

func TestNamespace_GetTable(t *testing.T) {
	n := NewNamespace(0)
	a := n.Get("abc", 1)
	if a == nil {
		t.Error("unexpected nil Table pointer")
	}
	b := n.Get("efg", 1)
	if b == nil {
		t.Error("unexpected nil Table pointer")
	}
	if a == b {
		t.Error("get different table name")
	}
	if a != n.Get("abc", 1) {
		t.Error("get the same table name")
	}
	a.Incr([]byte("11111"), 1)
	b.Incr([]byte("11111"), 1)
}

func TestNamespace_Destroy(t *testing.T) {
	n := NewNamespace(0)
	n.Get("a", 1)
	n.Get("b", 1)
	n.Get("c", 1)
	n.Destroy()
	if n.Data != nil {
		t.Errorf("expected namespace.Data nil pointer")
	}
}
