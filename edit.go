// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package edit implements buffered position-based editing of byte slices.
package edit

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// A Buffer is a queue of edits to apply to a given byte slice.
type Buffer struct {
	old []byte
	q   edits
}

// An edit records a single text modification: change the bytes in [start,end) to new.
type edit struct {
	start int
	end   int
	new   string
}

// An edits is a list of edits that is sortable by start offset, breaking ties by end offset.
type edits []edit

func (x edits) Len() int      { return len(x) }
func (x edits) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x edits) Less(i, j int) bool {
	if x[i].start != x[j].start {
		return x[i].start < x[j].start
	}
	return x[i].end < x[j].end
}

// NewBuffer returns a new buffer to accumulate changes to an initial data slice.
// The returned buffer maintains a reference to the data, so the caller must ensure
// the data is not modified until after the Buffer is done being used.
func NewBuffer(old []byte) *Buffer {
	return &Buffer{old: old}
}

// Insert inserts the new string at old[pos:pos].
func (b *Buffer) Insert(pos int, new string) {
	if pos < 0 || pos > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{pos, pos, new})
}

// Delete deletes the text old[start:end].
func (b *Buffer) Delete(start, end int) {
	if end < start || start < 0 || end > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{start, end, ""})
}

// Replace replaces old[start:end] with new.
func (b *Buffer) Replace(start, end int, new string) {
	if end < start || start < 0 || end > len(b.old) {
		panic("invalid edit position")
	}
	b.q = append(b.q, edit{start, end, new})
}

// Bytes returns a new byte slice containing the original data
// with the queued edits applied.
func (b *Buffer) Bytes() []byte {
	buf := new(bytes.Buffer)
	b.WriteTo(buf)
	return buf.Bytes()
}

// String returns a string containing the original data
// with the queued edits applied.
func (b *Buffer) String() string {
	buf := new(strings.Builder)
	b.WriteTo(buf)
	return buf.String()
}

// WriteTo writed the data with queued edits applied to w.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	// Sort edits by starting position and then by ending position.
	// Breaking ties by ending position allows insertions at point x
	// to be applied before a replacement of the text at [x, y).
	sort.Stable(b.q)

	var total int64
	write := func(p []byte) error {
		n, err := w.Write(p)
		total += int64(n)
		return err
	}
	writeStr := func(s string) error {
		n, err := io.WriteString(w, s)
		total += int64(n)
		return err
	}

	offset := 0
	for i, e := range b.q {
		if e.start < offset {
			e0 := b.q[i-1]
			panic(fmt.Sprintf("overlapping edits: [%d,%d)->%q, [%d,%d)->%q", e0.start, e0.end, e0.new, e.start, e.end, e.new))
		}
		if err := write(b.old[offset:e.start]); err != nil {
			return total, err
		}
		offset = e.end
		if err := writeStr(e.new); err != nil {
			return total, err
		}
	}
	err = write(b.old[offset:])
	return total, err
}
