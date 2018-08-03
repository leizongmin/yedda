package core

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable("")
	if table == nil {
		t.Error("expected Table pointer")
	}
}

func TestTable_Incr(t *testing.T) {
	table := NewTable("")
	if c := table.Incr([]byte("abc"), 1); c != 1 {
		t.Error("expected 1")
	}
	if c := table.Incr([]byte("abc"), 1); c != 2 {
		t.Error("expected 2")
	}
	if c := table.Incr([]byte("abc"), 1); c != 3 {
		t.Error("expected 3")
	}
	if c := table.Incr([]byte("123"), 2); c != 2 {
		t.Error("expected 2")
	}
	if c := table.Incr([]byte("123"), 3); c != 5 {
		t.Error("expected 5")
	}
}

func TestTable_Destroy(t *testing.T) {
	table := NewTable("")
	table.Incr([]byte("666"), 1)
	table.Incr([]byte("abcdefg"), 1)
	table.Destroy()
	if table.Hash != nil {
		t.Errorf("expected table.Hash nil pointer")
	}
	if table.Data != nil {
		t.Errorf("expected table.Data nil pointer")
	}
}
