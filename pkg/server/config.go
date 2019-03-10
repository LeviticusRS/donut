package server

import (
    "github.com/sprinkle-it/donut/pkg/client"
    "go.uber.org/zap"
)

type Config struct {
    ClientCapacity int
    LoggerConfig   zap.Config
    ClientConfig   client.Config
    Receivers      []client.MailReceiver
}

func (cfg Config) Build() (*Server, error) {
    logger, err := cfg.LoggerConfig.Build()
    if err != nil {
        return nil, err
    }

    router, err := client.NewMailRouter(cfg.Receivers)
    if err != nil {
        return nil, err
    }

    return &Server{
        logger:         logger,
        clientCapacity: cfg.ClientCapacity,
        clientFactory:  cfg.ClientConfig.Build,
        clients:        make(map[uint64]*client.Client, cfg.ClientCapacity),
        router:         router,
        commands:       make(chan command),
    }, nil
}
