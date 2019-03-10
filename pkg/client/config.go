package client

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
    "go.uber.org/zap"
    "net"
    "sync"
)

type Config struct {
    GenerateIdentifier IdentifierGenerator
    InputCapacity      int
    OutputCapacity     int
    MessageCapacity    int
}

func NewDefaultConfig() Config {
    return Config{
        GenerateIdentifier: IncrementalGenerator(0),
        InputCapacity:      10240,
        OutputCapacity:     10240,
        MessageCapacity:    1000,
    }
}

// Builds a new client from the configuration and given connection and router.
func (c *Config) Build(connection net.Conn, logger *zap.Logger, router MailRouter) *Client {
    return &Client{
        id:             c.GenerateIdentifier(),
        connection:     connection,
        logger:         logger,
        input:          buffer.NewRingBuffer(c.InputCapacity),
        output:         buffer.NewRingBuffer(c.OutputCapacity),
        outputCommands: make(chan outputCommand),
        decoder:        message.NewStreamDecoder(router.accepted, c.InputCapacity),
        encoder:        message.NewStreamEncoder(c.OutputCapacity),
        messages:       make(chan message.Message, c.MessageCapacity),
        router:         router,
        mutex:          sync.Mutex{},
        quit:           make(chan struct{}),
    }
}
