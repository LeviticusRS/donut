package player

import "github.com/sprinkle-it/donut/pkg/buffer"

type Walk struct {
    Direction uint8
}

func (w Walk) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(1, 2)
    buf.PutBits(uint32(w.Direction), 3)
    return nil
}

type Run struct {
    Direction uint8
}

func (r Run) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(1, 2)
    buf.PutBits(uint32(r.Direction), 4)
    return nil
}

type Teleported struct {

}

func (Teleported) Encode(*buffer.ByteBuffer) error {
    panic("implement me")
}

var Updated = updated{}

type updated struct {

}

func (u updated) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(1, 1)
    return nil
}


var NotUpdated = notUpdated{}

type notUpdated struct {

}

func (u notUpdated) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(0, 1)
    return nil
}

var Remove = remove{}

type remove struct {

}

func (remove) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(0, 1)
    buf.PutBits(0, 2)
    return nil
}