package fileservice

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    PassiveRequestDescriptor = message.Config{
        Id:   0,
        Size: 3,
        New:  func() message.Message { return &PassiveRequest{} },
    }

    PriorityRequestDescriptor = message.Config{
        Id:   1,
        Size: 3,
        New:  func() message.Message { return &PriorityRequest{} },
    }

    OnlineStatusUpdateDescriptor = message.Config{
        Id:   2,
        Size: 3,
        New:  message.Singleton(OnlineStatusUpdate),
    }

    OfflineStatusUpdateDescriptor = message.Config{
        Id:   3,
        Size: 3,
        New:  message.Singleton(OfflineStatusUpdate),
    }

    HandshakeDescriptor = message.Config{
        Id:   15,
        Size: 4,
        New:  func() message.Message { return &Handshake{} },
    }
)

type Handshake struct {
    Version uint32
}

func (Handshake) Config() message.Config { return HandshakeDescriptor }

func (h *Handshake) Decode(buf *buffer.ByteBuffer, length int) error {
    var err error
    if h.Version, err = buf.GetUint32(); err != nil {
        return err
    }
    return nil
}

type Request struct {
    Index uint8
    Id    uint16
}

type PassiveRequest struct{ Request }

func (r PassiveRequest) Config() message.Config {
    return PassiveRequestDescriptor
}

func (r *PassiveRequest) Decode(buf *buffer.ByteBuffer, length int) error {
    var err error

    if r.Index, err = buf.GetUint8(); err != nil {
        return err
    }

    if r.Id, err = buf.GetUint16(); err != nil {
        return err
    }

    return nil
}

type PriorityRequest struct{ Request }

func (r PriorityRequest) Config() message.Config {
    return PriorityRequestDescriptor
}

func (r *PriorityRequest) Decode(buf *buffer.ByteBuffer, length int) error {
    var err error

    if r.Index, err = buf.GetUint8(); err != nil {
        return err
    }

    if r.Id, err = buf.GetUint16(); err != nil {
        return err
    }

    return nil
}

var OnlineStatusUpdate = onlineStatusUpdate{}

type onlineStatusUpdate struct{}

func (s onlineStatusUpdate) Config() message.Config { return OnlineStatusUpdateDescriptor }

func (onlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

var OfflineStatusUpdate = offlineStatusUpdate{}

type offlineStatusUpdate struct{}

func (offlineStatusUpdate) Config() message.Config { return OfflineStatusUpdateDescriptor }

func (offlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error { return nil }