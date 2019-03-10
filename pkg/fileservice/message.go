package fileservice

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    PassiveRequestDescriptor      = message.Descriptor{Id: 0, Size: 3, Provider: newPassiveRequest,}
    PriorityRequestDescriptor     = message.Descriptor{Id: 1, Size: 3, Provider: newPriorityRequest,}
    OnlineStatusUpdateDescriptor  = message.Descriptor{Id: 2, Size: 3, Provider: message.Provide(OnlineStatusUpdate),}
    OfflineStatusUpdateDescriptor = message.Descriptor{Id: 3, Size: 3, Provider: message.Provide(OfflineStatusUpdate),}
    HandshakeDescriptor           = message.Descriptor{Id: 15, Size: 4, Provider: newHandshake,}
)

type Handshake struct {
    Version uint32
}

func newHandshake() message.Message { return &Handshake{} }

func (Handshake) Descriptor() message.Descriptor { return HandshakeDescriptor }

func (h *Handshake) Decode(buf *buffer.ByteBuffer, length int) error {
    var err error
    if h.Version, err = buf.GetUint32(); err != nil {
        return err
    }
    return nil
}

func (Handshake) Encode(*buffer.ByteBuffer) error { return message.ErrEncodeNotSupported }

type Request struct {
    Index uint8
    Id    uint16
}

type PassiveRequest struct{ Request }

func newPassiveRequest() message.Message { return &PassiveRequest{} }

func (r PassiveRequest) Descriptor() message.Descriptor {
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

func (PassiveRequest) Encode(*buffer.ByteBuffer) error { return message.ErrEncodeNotSupported }

type PriorityRequest struct{ Request }

func newPriorityRequest() message.Message { return &PriorityRequest{} }

func (r PriorityRequest) Descriptor() message.Descriptor {
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

func (PriorityRequest) Encode(*buffer.ByteBuffer) error { return message.ErrEncodeNotSupported }

var OnlineStatusUpdate = onlineStatusUpdate{}

type onlineStatusUpdate struct{}

func (s onlineStatusUpdate) Descriptor() message.Descriptor { return OnlineStatusUpdateDescriptor }

func (onlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

func (onlineStatusUpdate) Encode(*buffer.ByteBuffer) error { return message.ErrEncodeNotSupported }

var OfflineStatusUpdate = offlineStatusUpdate{}

type offlineStatusUpdate struct{}

func (offlineStatusUpdate) Descriptor() message.Descriptor { return OfflineStatusUpdateDescriptor }

func (offlineStatusUpdate) Decode(buf *buffer.ByteBuffer, length int) error { return nil }

func (offlineStatusUpdate) Encode(*buffer.ByteBuffer) error { return message.ErrEncodeNotSupported }
