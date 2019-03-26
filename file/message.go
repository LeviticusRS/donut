package file

import (
    "github.com/sprinkle-it/donut/buffer"
    "github.com/sprinkle-it/donut/message"
)

var (
    passiveRequestConfig = message.Config{
        Id:   0,
        Size: 3,
        New:  func() message.Message { return &PassiveRequest{} },
    }

    priorityRequestConfig = message.Config{
        Id:   1,
        Size: 3,
        New:  func() message.Message { return &PriorityRequest{} },
    }

    onlineStatusUpdateConfig = message.Config{
        Id:   2,
        Size: 3,
        New:  message.Singleton(OnlineStatusUpdate),
    }

    offlineStatusUpdateConfig = message.Config{
        Id:   3,
        Size: 3,
        New:  message.Singleton(OfflineStatusUpdate),
    }

    handshakeConfig = message.Config{
        Id:   15,
        Size: 4,
        New:  func() message.Message { return &Handshake{} },
    }
)

type Handshake struct {
    Version uint32
}

func (Handshake) Config() message.Config { return handshakeConfig }

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

func (r PassiveRequest) Config() message.Config { return passiveRequestConfig }

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
    return priorityRequestConfig
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

func (s onlineStatusUpdate) Config() message.Config { return onlineStatusUpdateConfig }

func (onlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error { return nil }

var OfflineStatusUpdate = offlineStatusUpdate{}

type offlineStatusUpdate struct{}

func (offlineStatusUpdate) Config() message.Config { return offlineStatusUpdateConfig }

func (offlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error { return nil }