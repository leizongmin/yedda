package core

import "testing"

func TestNewDataBase(t *testing.T) {
	d := NewDataBase(16)
	if d == nil {
		t.Error("unexpected nil Database pointer")
	}
	if len(d.Data) != 16 {
		t.Error("unexpected database size")
	}
	d2 := NewDataBase(32)
	if d2 == nil {
		t.Error("unexpected nil Database pointer")
	}
	if len(d2.Data) != 32 {
		t.Error("unexpected database size")
	}
	d3 := NewDataBase(512)
	if d3 == nil {
		t.Error("unexpected nil Database pointer")
	}
	if len(d3.Data) != 512 {
		t.Error("unexpected database size")
	}

}
