// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package edit

import "testing"

func TestEdit(t *testing.T) {
	b := NewBuffer([]byte("0123456789"))
	b.Insert(8, ",7½,")
	b.Replace(9, 10, "the-end")
	b.Insert(10, "!")
	b.Insert(4, "3.14,")
	b.Insert(4, "π,")
	b.Insert(4, "3.15,")
	b.Replace(3, 4, "three,")
	want := "012three,3.14,π,3.15,4567,7½,8the-end!"

	s := b.String()
	if s != want {
		t.Errorf("b.String() = %q, want %q", s, want)
	}
	sb := b.Bytes()
	if string(sb) != want {
		t.Errorf("b.Bytes() = %q, want %q", sb, want)
	}
}

func TestOverlappingDeletes(t *testing.T) {
	const in = "0123456789"
	const want = "0156789"

	// Test single delete.
	b := NewBuffer([]byte(in))
	b.Delete(2, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}

	// Test overlap at beginning.
	b = NewBuffer([]byte(in))
	b.Delete(2, 3)
	b.Delete(2, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}

	// Test overlap in middle.
	b = NewBuffer([]byte(in))
	b.Delete(3, 4)
	b.Delete(2, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}

	// Test overlap at end.
	b = NewBuffer([]byte(in))
	b.Delete(4, 5)
	b.Delete(2, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}

	// Test covering overlap.
	b = NewBuffer([]byte(in))
	b.Delete(2, 3)
	b.Delete(3, 5)
	b.Delete(2, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}

	// Test partial overlap.
	b = NewBuffer([]byte(in))
	b.Delete(2, 4)
	b.Delete(3, 5)
	if got := b.String(); got != want {
		t.Errorf("b.String() = %q want %q", got, want)
	}
}

var sink []byte

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b := NewBuffer([]byte("0123456789"))
		b.Insert(8, ",7½,")
		b.Replace(9, 10, "the-end")
		b.Insert(10, "!")
		b.Insert(4, "3.14,")
		b.Insert(4, "π,")
		b.Insert(4, "3.15,")
		b.Replace(3, 4, "three,")
		sink = b.Bytes()
	}
}
