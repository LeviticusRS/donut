package buffer

import "errors"

type ByteBuffer struct {
    Bytes  []byte
    Offset int
}

func NewByteBuffer(capacity int) ByteBuffer {
    return ByteBuffer{Bytes: make([]byte, capacity)}
}

func (b *ByteBuffer) check(n int) error {
    if len(b.Bytes)-b.Offset < n {
        return errors.New("buffer: insufficient bytes remaining to perform this operation")
    }
    return nil
}

func (b *ByteBuffer) Skip(amount int) error {
    if err := b.check(amount); err != nil {
        return err
    }
    b.Offset += amount
    return nil
}

func (b *ByteBuffer) GetUint8() (uint8, error) {
    if err := b.check(1); err != nil {
        return 0, err
    }
    b.Offset += 1
    v := b.Bytes[b.Offset-1]
    return v, nil
}

func (b *ByteBuffer) PutUint8(v uint8) error {
    if err := b.check(1); err != nil {
        return err
    }
    b.Offset += 1
    b.Bytes[b.Offset-1] = uint8(v)
    return nil
}

func (b *ByteBuffer) GetUint16() (uint16, error) {
    if err := b.check(2); err != nil {
        return 0, err
    }
    b.Offset += 2
    v := uint16(b.Bytes[b.Offset-2])<<8 | uint16(b.Bytes[b.Offset-1])
    return v, nil
}

func (b *ByteBuffer) PutUint16(v uint16) error {
    if err := b.check(2); err != nil {
        return err
    }
    b.Offset += 2
    b.Bytes[b.Offset-2] = uint8(v >> 8)
    b.Bytes[b.Offset-1] = uint8(v)
    return nil
}

func (b *ByteBuffer) GetUint32() (uint32, error) {
    if err := b.check(4); err != nil {
        return 0, err
    }
    b.Offset += 4
    v := uint32(b.Bytes[b.Offset-4])<<24 | uint32(b.Bytes[b.Offset-3])<<16 |
        uint32(b.Bytes[b.Offset-2])<<8 | uint32(b.Bytes[b.Offset-1])
    return v, nil
}

func (b *ByteBuffer) PutUint32(v uint32) error {
    if err := b.check(4); err != nil {
        return err
    }

    b.Offset += 4
    b.Bytes[b.Offset-4] = uint8(v >> 24)
    b.Bytes[b.Offset-3] = uint8(v >> 16)
    b.Bytes[b.Offset-2] = uint8(v >> 8)
    b.Bytes[b.Offset-1] = uint8(v)
    return nil
}

func (b *ByteBuffer) GetUint64() (uint64, error) {
    if err := b.check(8); err != nil {
        return 0, err
    }
    b.Offset += 8
    v := uint64(b.Bytes[b.Offset-8])<<56 | uint64(b.Bytes[b.Offset-7])<<48 |
        uint64(b.Bytes[b.Offset-6])<<40 | uint64(b.Bytes[b.Offset-5])<<32 |
        uint64(b.Bytes[b.Offset-4])<<24 | uint64(b.Bytes[b.Offset-3])<<16 |
        uint64(b.Bytes[b.Offset-2])<<8 | uint64(b.Bytes[b.Offset-1])
    return v, nil
}

func (b *ByteBuffer) PutUint64(v uint64) error {
    if err := b.check(8); err != nil {
        return err
    }
    b.Offset += 4
    b.Bytes[b.Offset-8] = uint8(v >> 56)
    b.Bytes[b.Offset-7] = uint8(v >> 48)
    b.Bytes[b.Offset-6] = uint8(v >> 40)
    b.Bytes[b.Offset-5] = uint8(v >> 32)
    b.Bytes[b.Offset-4] = uint8(v >> 24)
    b.Bytes[b.Offset-3] = uint8(v >> 16)
    b.Bytes[b.Offset-2] = uint8(v >> 8)
    b.Bytes[b.Offset-1] = uint8(v)
    return nil
}
