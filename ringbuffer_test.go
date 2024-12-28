package ringbuffer_test

import (
	"fmt"
	"github.com/nsf/ringbuffer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	assert := assert.New(t)

	eq2 := func(v int, ok bool) func(expectedV int, expectedOk bool) {
		return func(expectedV int, expectedOk bool) {
			assert.Equal(expectedOk, ok)
			assert.Equal(expectedV, v)
		}
	}

	{
		var buf ringbuffer.RingBuffer[int]
		assert.Equal(false, buf.Push(5))
		assert.Equal(0, buf.Cap())
		assert.Equal(0, buf.Len())
		eq2(buf.Pop())(0, false)
	}

	{
		buf := ringbuffer.New[int](1)
		eq2(buf.Pop())(0, false)
		assert.Equal(1, buf.Cap())
		assert.Equal(0, buf.Len())

		assert.Equal(true, buf.Push(3))
		assert.Equal(1, buf.Cap())
		assert.Equal(1, buf.Len())

		assert.Equal(false, buf.Push(3))
		assert.Equal(1, buf.Cap())
		assert.Equal(1, buf.Len())

		eq2(buf.Pop())(3, true)
		assert.Equal(1, buf.Cap())
		assert.Equal(0, buf.Len())
		eq2(buf.Pop())(0, false)
		assert.Equal(1, buf.Cap())
		assert.Equal(0, buf.Len())
	}

	buf := ringbuffer.New[int](3)
	for i := 0; i < 10; i++ {
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())

		assert.Equal(true, buf.Push(1))
		assert.Equal(3, buf.Cap())
		assert.Equal(1, buf.Len())

		assert.Equal(true, buf.Push(2))
		assert.Equal(3, buf.Cap())
		assert.Equal(2, buf.Len())

		assert.Equal(true, buf.Push(3))
		assert.Equal(3, buf.Cap())
		assert.Equal(3, buf.Len())

		eq2(buf.Pop())(1, true)
		assert.Equal(3, buf.Cap())
		assert.Equal(2, buf.Len())

		eq2(buf.Pop())(2, true)
		assert.Equal(3, buf.Cap())
		assert.Equal(1, buf.Len())

		eq2(buf.Pop())(3, true)
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())

		eq2(buf.Pop())(0, false)
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())
	}
	for i := 0; i < 10; i++ {
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())

		assert.Equal(true, buf.Push(1))
		assert.Equal(3, buf.Cap())
		assert.Equal(1, buf.Len())

		assert.Equal(true, buf.Push(2))
		assert.Equal(3, buf.Cap())
		assert.Equal(2, buf.Len())

		eq2(buf.Pop())(1, true)
		assert.Equal(3, buf.Cap())
		assert.Equal(1, buf.Len())

		eq2(buf.Pop())(2, true)
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())

		eq2(buf.Pop())(0, false)
		assert.Equal(3, buf.Cap())
		assert.Equal(0, buf.Len())

		assert.Equal(true, buf.Push(1))
		assert.Equal(true, buf.Push(2))
		assert.Equal(true, buf.Push(3))
		assert.Equal(false, buf.Push(4))
		eq2(buf.Pop())(1, true)
		eq2(buf.Pop())(2, true)
		eq2(buf.Pop())(3, true)
		eq2(buf.Pop())(0, false)
	}
}

func ExampleRingBuffer() {
	// Using ringbuffer structure alone without the "New" function is fairly useless, but it's valid.
	var buf ringbuffer.RingBuffer[int]
	fmt.Printf("cap: %d, len: %d, pushed: %v\n", buf.Cap(), buf.Len(), buf.Push(5))
	// Output: cap: 0, len: 0, pushed: false
}

func ExampleNew() {
	b := ringbuffer.New[int](5)
	b.Push(1)
	b.Push(2)
	v1, _ := b.Pop()
	v2, _ := b.Pop()
	fmt.Printf("%d %d\n", v1, v2)
	// Output: 1 2
}

func ExampleRingBuffer_Cap() {
	b := ringbuffer.New[int](5)
	fmt.Printf("%d\n", b.Cap())
	// Output: 5
}

func ExampleRingBuffer_Len() {
	b := ringbuffer.New[int](5)
	l1 := b.Len()
	b.Push(1)
	l2 := b.Len()
	fmt.Printf("%d %d\n", l1, l2)
	// Output: 0 1
}

func ExampleRingBuffer_Push() {
	b := ringbuffer.New[int](5)
	b.Push(1)
	b.Push(2)
	v1, _ := b.Pop()
	v2, _ := b.Pop()
	fmt.Printf("%d %d\n", v1, v2)
	// Output: 1 2
}

func ExampleRingBuffer_Pop() {
	b := ringbuffer.New[int](5)
	b.Push(1)
	b.Push(2)
	v1, _ := b.Pop()
	v2, _ := b.Pop()
	fmt.Printf("%d %d\n", v1, v2)
	// Output: 1 2
}

func ExamplePush() {
	var buf [5]int
	var read int8
	var write int8

	ringbuffer.Push(buf[:], read, &write, 1)
	ringbuffer.Push(buf[:], read, &write, 2)
	v1, _ := ringbuffer.Pop(buf[:], &read, write)
	v2, _ := ringbuffer.Pop(buf[:], &read, write)
	fmt.Printf("%d %d\n", v1, v2)
	// Output: 1 2
}

func ExamplePop() {
	var buf [5]int
	var read int8
	var write int8

	ringbuffer.Push(buf[:], read, &write, 1)
	ringbuffer.Push(buf[:], read, &write, 2)
	v1, _ := ringbuffer.Pop(buf[:], &read, write)
	v2, _ := ringbuffer.Pop(buf[:], &read, write)
	fmt.Printf("%d %d\n", v1, v2)
	// Output: 1 2
}

func ExampleCap() {
	var buf [10]int
	fmt.Printf("%d\n", ringbuffer.Cap(buf[:]))
	// Output: 9
}

func ExampleLen() {
	var buf [10]int
	var read int8
	var write int8

	l1 := ringbuffer.Len(buf[:], read, write)
	ringbuffer.Push(buf[:], read, &write, 1)
	l2 := ringbuffer.Len(buf[:], read, write)
	fmt.Printf("%d %d\n", l1, l2)
	// Output: 0 1
}
