package game

import "github.com/sprinkle-it/donut/pkg/buffer"

type Position struct {
    Level uint8
    X     uint16
    Z     uint16
}

func (p Position) EncodeHash(buf *buffer.ByteBuffer) {
    buf.WriteBits(uint32(p.Level), 2)
    buf.WriteBits(uint32(p.X), 14)
    buf.WriteBits(uint32(p.Z), 14)
}

func (p Position) EncodeBlockHash(buf *buffer.ByteBuffer) {
    buf.WriteBits(uint32(p.Level), 2)
    buf.WriteBits(uint32(p.X >> 6), 8)
    buf.WriteBits(uint32(p.Z >> 6), 8)
}