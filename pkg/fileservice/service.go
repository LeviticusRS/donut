package fileservice

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/client"
    "github.com/sprinkle-it/donut/pkg/status"
    "go.uber.org/zap"
)

type ArchiveProvider func(uint8, uint16) ([]byte, error)

type Service struct {
    logger        *zap.Logger
    capacity      int
    sessions      map[uint64]*Session
    newSession    SessionFactory
    getArchive    ArchiveProvider
    commands      chan command
    workers       WorkerPool
    clientVersion uint32
}

func New(config Config) (*Service, error) {
    return config.Build()
}

func (s *Service) execute(cmd command) { s.commands <- cmd }

func (s *Service) Process() {
    for i := 0; i < cap(s.workers); i++ {
        worker := Worker{
            pool:     s.workers,
            jobs:     make(JobQueue),
            provider: s.getArchive,
        }
        worker.Process()
    }

    go func() {
        for command := range s.commands {
            command.execute(s)
        }
    }()
}

func (s *Service) HandleMail(mail client.Mail) {
    s.execute(handleMessage{mail: mail})
}

type command interface {
    execute(s *Service)
}

type handleMessage struct {
    mail client.Mail
}

func (c handleMessage) execute(s *Service) {
    source := c.mail.Source
    switch msg := c.mail.Message.(type) {
    case *Handshake:
        if msg.Version != s.clientVersion {
            _ = source.SendNow(status.UnsupportedVersion)
            return
        }

        if len(s.sessions) >= s.capacity {
            _ = source.SendNow(status.Full)
            return
        }

        session := s.newSession(source, s.workers)
        s.sessions[source.Id()] = session

        source.OnClosed(func(cli *client.Client) {
            s.execute(unregisterSession{client: cli})
        })

        session.Process()

        session.Info("Registered client to file service")

        _ = source.SendNow(status.Okay)
    case *PassiveRequest:
        session, exists := s.sessions[source.Id()]
        if !exists {
            source.Fatal(errors.New("fileservice: received request from client that does not have active session"))
            return
        }
        session.enqueue(msg.Request, session.passive)
    case *PriorityRequest:
        session, exists := s.sessions[source.Id()]
        if !exists {
            source.Fatal(errors.New("fileservice: received request from client that does not have active session"))
            return
        }
        session.enqueue(msg.Request, session.priority)
    }
}

type unregisterSession struct {
    client *client.Client
}

func (cmd unregisterSession) execute(service *Service) {
    delete(service.sessions, cmd.client.Id())

    service.logger.Info("Unregistered file session",
        zap.Uint64("id", cmd.client.Id()),
        zap.Stringer("address", cmd.client.RemoteAddress()),
    )
}
