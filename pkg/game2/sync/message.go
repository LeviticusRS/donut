package sync

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    messageConfig = message.Config{
        Id:   79,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &Message{} },
    }
)

type Block interface {
    Encode(*buffer.ByteBuffer) error
}

var StartList = startList{}

var EndList = endList{}

type startList struct{}

func (startList) Encode(buf *buffer.ByteBuffer) error {
    buf.StartBitAccess()
    return nil
}

type endList struct{}

func (endList) Encode(buf *buffer.ByteBuffer) error {
    buf.FinishBitAccess()
    return nil
}

type Skip struct {
    Count int
}

func (b Skip) Encode(buf *buffer.ByteBuffer) error {
    // Flag for if the block is a descriptor, skips blocks are not descriptors.
    buf.PutBits(0, 1)

    // Write out the amount of players in the list to skip.
    switch {
    case b.Count < 1:
        buf.PutBits(0, 2)
        return nil
    case b.Count >= 1 && b.Count <= 31:
        buf.PutBits(1, 2)
        buf.PutBits(uint32(b.Count), 5)
    case b.Count >= 32 && b.Count <= 255:
        buf.PutBits(2, 2)
        buf.PutBits(uint32(b.Count), 8)
    default:
        buf.PutBits(3, 2)
        buf.PutBits(uint32(b.Count), 11)
    }
    return nil
}

type BlockList []Block

func (l *BlockList) Push(block Block) {
    *l = append(*l, block)
}

func (l *BlockList) Copy() BlockList {
    list := make(BlockList, len(*l))
    copy(list, *l)
    return list
}

func (l *BlockList) Clear() {
    old := *l
    *l = old[:0]
}

type Message struct {
    Blocks BlockList
}

func (Message) Config() message.Config { return messageConfig }

func (u Message) Encode(buf *buffer.ByteBuffer) error {
    for _, block := range u.Blocks {
        if err := block.Encode(buf); err != nil {
            return err
        }
    }
    return nil
}