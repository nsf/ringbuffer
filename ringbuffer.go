// Package provides fixed length FIFO ring buffer functionality.
//
// You can use a generic data structure with methods, which is forced to use slices due to Go language limitations
// (no generic array lengths):
//
//	b := ringbuffer.New[int](5)
//	b.Push(1)
//	b.Push(2)
//	v1, _ := b.Pop()
//	v2, _ := b.Pop()
//	fmt.Printf("%d %d\n", v1, v2)
//
// Or you can use generic functions with your own state variables. Here you can use plain arrays and integer types of
// arbitrary size:
//
//	var buf [5]int
//	var read int8
//	var write int8
//
//	ringbuffer.Push(buf[:], read, &write, 1)
//	ringbuffer.Push(buf[:], read, &write, 2)
//	v1, _ := ringbuffer.Pop(buf[:], &read, write)
//	v2, _ := ringbuffer.Pop(buf[:], &read, write)
//	fmt.Printf("%d %d\n", v1, v2)
//
// The logic is the simplest implementation straight from wikipedia: https://en.wikipedia.org/wiki/Circular_buffer. In
// short: there are read and write pointers as integer indices and a buffer of capacity+1 space. An extra element is
// reserved to distinguish between full/empty state. When using plain generic functions with smaller integer types, be
// aware of integer overflows. No care taken to prevent these.
package ringbuffer

import (
	"golang.org/x/exp/constraints"
)

// Fixed length FIFO ring buffer.
type RingBuffer[T any] struct {
	read   int
	write  int
	buffer []T
}

// Create a new buffer which can store capacity elements. The buffer is fixed in length and will not grow.
// It is a FIFO buffer.
//
// Implementation detail: the buffer is implemented as two integer pointers and a slice of capacity+1 elements.
// One extra element is reserved to avoid ambiguous state where read and write pointers point to the same location
// and it might mean full or empty buffer. A bit of extra space wasted is traded for logical simplicity.
func New[T any](capacity int) RingBuffer[T] {
	var buffer []T
	if capacity >= 1 {
		buffer = make([]T, capacity+1)
	}
	return RingBuffer[T]{
		buffer: buffer,
		read:   0,
		write:  0,
	}
}

// How many elements a buffer can store?
func (b RingBuffer[T]) Cap() int {
	return Cap(b.buffer)
}

// How many elements are currently stored in the buffer?
func (b RingBuffer[T]) Len() int {
	return Len(b.buffer, b.read, b.write)
}

// Push a new element to the buffer.
//
// Returns true on success. Returns false if there is no free space and push failed.
func (b *RingBuffer[T]) Push(v T) bool {
	return Push(b.buffer, b.read, &b.write, v)
}

// Try to pop an element from the buffer.
//
// Returns the popped element and true on success. Returns default value and false if there were no elements in the buffer.
func (b *RingBuffer[T]) Pop() (T, bool) {
	return Pop(b.buffer, &b.read, b.write)
}

// How many elements a buffer can store?
func Cap[T any](slice []T) int {
	v := len(slice) - 1
	if v < 0 {
		v = 0
	}
	return v
}

// How many elements are currently stored in the buffer?
func Len[T any, U constraints.Integer](slice []T, read, write U) int {
	if write >= read {
		return int(write - read)
	} else {
		return len(slice) - int(read-write)
	}
}

// Push a new element to the buffer.
//
// Returns true on success. Returns false if there is no free space and push failed.
func Push[T any, U constraints.Integer](slice []T, read U, write *U, v T) bool {
	if len(slice) == 0 {
		return false
	}
	next := int(*write+1) % len(slice)
	if next == int(read) {
		return false // no more space
	}
	slice[*write] = v
	*write = U(next)
	return true
}

// Try to pop an element from the buffer.
//
// Returns the popped element and true on success. Returns default value and false if there were no elements in the buffer.
func Pop[T any, U constraints.Integer](slice []T, read *U, write U) (T, bool) {
	if *read == write {
		var def T
		return def, false
	}
	val := slice[*read]
	*read = U(int(*read+1) % len(slice))
	return val, true
}
