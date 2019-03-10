package fileservice

import (
    "go.uber.org/zap"
)

type Config struct {
    LoggerConfig     zap.Config
    Capacity         int
    Workers          int
    SupportedVersion uint32
    Provider         ArchiveProvider
    SessionConfig    SessionConfig
}

func (cfg Config) Build() (*Service, error) {
    logger, err := cfg.LoggerConfig.Build()
    if err != nil {
        return nil, err
    }

    return &Service{
        logger:        logger,
        capacity:      cfg.Capacity,
        commands:      make(chan command),
        sessions:      make(map[uint64]*Session, cfg.Capacity),
        newSession:    cfg.SessionConfig.Build,
        clientVersion: cfg.SupportedVersion,
        workers:       make(WorkerQueue, cfg.Workers),
        getArchive:    cfg.Provider,
    }, nil
}
