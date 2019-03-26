package game

import (
    "github.com/sprinkle-it/donut/buffer"
    "github.com/sprinkle-it/donut/message"
)

var (
    handshakeConfig = message.Config{
        Id:   14,
        Size: 0,
        New:  message.Singleton(Handshake),
    }

    authenticateConfig = message.Config{
        Id:   16,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &Authenticate{} },
    }

    Handshake    = &handshake{}
)

type handshake struct{}

func (handshake) Config() message.Config { return handshakeConfig }

func (handshake) Decode(buf *buffer.ByteBuffer, length int) error { return nil }

type Ready struct {
    AuthenticationKey uint64
}

type Authenticate struct {
}

func (Authenticate) Config() message.Config { return authenticateConfig }

func (Authenticate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}
