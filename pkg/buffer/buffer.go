package buffer

import "io"

// Type declaration for readable buffers. Readable buffers must implement io.Reader and specify
// how many readable bytes they have remaining.
type Readable interface {
    io.Reader
    Readable() int
}

// Checks if the readable buffer has any available readable bytes.
func HasReadable(r Readable) bool {
    return r.Readable() > 0
}

// Checks if the readable buffer has a specified amount of readable bytes.
func IsReadable(r Readable, n int) bool {
    return r.Readable() >= n
}

type Writable interface {
    io.Writer
    Writable() int
}

func HasWritable(w Writable) bool {
    return w.Writable() > 0
}