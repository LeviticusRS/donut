package game2

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    HandshakeConfig = message.Config{
        Id:   14,
        Size: 0,
        New:  message.Singleton(Handshake),
    }

    AuthenticateConfig = message.Config{
        Id:   16,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &Authenticate{} },
    }

    Handshake    = &handshake{}
)

type handshake struct{}

func (handshake) Config() message.Config { return HandshakeConfig }

func (handshake) Decode(buf *buffer.ByteBuffer, length int) error { return nil }

type Ready struct {
    AuthenticationKey uint64
}

type Authenticate struct {
}

func (Authenticate) Config() message.Config {
    return AuthenticateConfig
}

func (Authenticate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}
