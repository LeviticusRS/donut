package fileservice

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    PassiveRequestConfig = message.Config{
        Id:   0,
        Size: 3,
        New:  func() message.Message { return &PassiveRequest{} },
    }

    PriorityRequestConfig = message.Config{
        Id:   1,
        Size: 3,
        New:  func() message.Message { return &PriorityRequest{} },
    }

    OnlineStatusUpdateConfig = message.Config{
        Id:   2,
        Size: 3,
        New:  message.Singleton(OnlineStatusUpdate),
    }

    OfflineStatusUpdateConfig = message.Config{
        Id:   3,
        Size: 3,
        New:  message.Singleton(OfflineStatusUpdate),
    }

    HandshakeConfig = message.Config{
        Id:   15,
        Size: 4,
        New:  func() message.Message { return &Handshake{} },
    }
)

type Handshake struct {
    Version uint32
}

func (Handshake) Config() message.Config { return HandshakeConfig }

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
    return PassiveRequestConfig
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
    return PriorityRequestConfig
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

func (s onlineStatusUpdate) Config() message.Config { return OnlineStatusUpdateConfig }

func (onlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

var OfflineStatusUpdate = offlineStatusUpdate{}

type offlineStatusUpdate struct{}

func (offlineStatusUpdate) Config() message.Config { return OfflineStatusUpdateConfig }

func (offlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error { return nil }