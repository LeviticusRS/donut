package game

import "github.com/sprinkle-it/donut/pkg/buffer"

type Position struct {
    Level uint8
    X     uint16
    Z     uint16
}

func (p Position) ChunkX() uint16 {
    return p.X >> 3
}

func (p Position) ChunkZ() uint16 {
    return p.Z >> 3
}

func (p Position) EncodeHash(buf *buffer.ByteBuffer) {
    buf.PutBits(uint32(p.Level), 2)
    buf.PutBits(uint32(p.X), 14)
    buf.PutBits(uint32(p.Z), 14)
}

func (p Position) EncodeBlockHash(buf *buffer.ByteBuffer) {
    buf.PutBits(uint32(p.Level), 2)
    buf.PutBits(uint32(p.X >> 6), 8)
    buf.PutBits(uint32(p.Z >> 6), 8)
}