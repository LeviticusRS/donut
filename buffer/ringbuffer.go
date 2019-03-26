package buffer

import (
    "errors"
    "io"
)

type RingBuffer struct {
    capacity int
    bytes    []byte
    readPos  int
    writePos int
}

func NewRingBuffer(capacity int) RingBuffer {
    return RingBuffer{
        capacity: capacity,
        bytes: make([]byte, capacity+1),
    }
}

func (b *RingBuffer) Capacity() int {
    return b.capacity
}

func (b *RingBuffer) Read(bytes []byte) (int, error) {
    if len(bytes) == 0 {
        return 0, nil
    }

    n := min(len(bytes), b.Readable())

    if b.readPos+n > len(b.bytes) {
        // First copy to the destination array the bytes between the read position and the end of the internal buffer.
        copy(bytes[:len(b.bytes)-b.readPos], b.bytes[b.readPos:])

        // Lastly copy to the destination array the remaining bytes at the start of the internal buffer.
        copy(bytes[len(b.bytes)-b.readPos:], b.bytes[:n-len(b.bytes)+b.readPos])
    } else {
        // Copy the requested bytes from the read position.
        copy(bytes, b.bytes[b.readPos:])
    }

    b.readPos = (b.readPos + n) % len(b.bytes)

    // Implementations of Read are discouraged from returning a zero byte count with a nil error.
    if n == 0 {
        return n, errors.New("")
    }

    return n, nil
}

func (b *RingBuffer) Write(bytes []byte) (int, error) {
    n := min(len(bytes), b.Writable())

    if b.writePos+n > len(b.bytes) {
        // First copy to the write position of the internal array the bytes at the start of the source buffer.
        copy(b.bytes[b.writePos:], bytes[:len(b.bytes)-b.writePos])

        // Lastly copy to the start of the internal array the remaining bytes after what we had just written.
        copy(b.bytes[:n-len(b.bytes)+b.writePos], bytes[len(b.bytes)-b.writePos:])
    } else {
        // Copy the bytes to the write position
        copy(b.bytes[b.writePos:], bytes)
    }

    b.writePos = (b.writePos + n) % len(b.bytes)

    // Write must return a non-nil error if it returns n < len(p).
    if n != len(bytes) {
        return n, io.ErrShortWrite
    }

    return n, nil
}

func (b *RingBuffer) Readable() int {
    // Check if the write position has not wrapped over the boundary to be before the read position.
    if b.readPos <= b.writePos {
        return b.writePos - b.readPos
    }
    return len(b.bytes) - b.readPos + b.writePos
}

func (b *RingBuffer) Writable() int {
    // Check that the write position has not wrapped over the boundary to be before the read position. Reserve
    // the last byte before the read position as a buffer so that we can know when it is full versus just when
    // it is initialized.
    if b.readPos <= b.writePos {
        return len(b.bytes) - b.writePos + b.readPos - 1
    }
    return b.readPos - b.writePos - 1
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
