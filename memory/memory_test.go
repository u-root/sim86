package mem

import (
	"testing"
)
func TestLong(t*testing.T) {
	memory = make([]uint32, 262144)
	r := NewLongMemReader(0)
	t.Log("r.Read is %v", r.Read)
	w := NewLongMemWriter(0)
	w.Write(0xdeadbeef)
	if r.Read() != 0xdeadbeef {
		t.Errorf("r: got %v, want %v", r.Read(), 0xdeadbeef)
	}
	rb := NewByteMemReader(0)
	if rb.Read() != 0xef {
		t.Errorf("r: got %v, want %v", rb.Read(), 0xef)
	}
	wb := NewByteMemWriter(0)
	wb.Write(0)
	if r.Read() != 0xdeadbe00 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0xdeadbe00)
	}
	wb = NewByteMemWriter(1)
	wb.Write(1)
	if r.Read() != 0xdead0100 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0xdead0100)
	}
	wb = NewByteMemWriter(2)
	wb.Write(2)
	if r.Read() != 0xde020100 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0xde020100)
	}
	wb = NewByteMemWriter(3)
	wb.Write(3)
	if r.Read() != 0x03020100 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0x03020100)
	}

	ww := NewWordMemWriter(0)
	ww.Write(0xaa55)
	if r.Read() != 0x0302aa55 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0x0302aa55)
	}
	ww = NewWordMemWriter(1)
	ww.Write(0xcafe)
	if r.Read() != 0xcafeaa55 {
		t.Errorf("r: got %#08x, want %#08x", r.Read(), 0xcafeaa55)
	}

}


