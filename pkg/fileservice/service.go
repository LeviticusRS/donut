package fileservice

import (
    "errors"
    "github.com/sprinkle-it/donut/pkg/client"
    "github.com/sprinkle-it/donut/pkg/status"
    "go.uber.org/zap"
)

type ArchiveProvider func(uint8, uint16) ([]byte, error)

type Service struct {
    logger *zap.Logger

    // The number of concurrent sessions the service can have up to. If a client attempts to connect to this service
    // and the service is at capacity the service will reply with a status of Full and close the client.
    capacity int

    // The sessions that are currently active for this service.
    sessions map[uint64]*Session

    // Factory function to create new sessions.
    newSession SessionFactory

    // The client version that the service supports. If a client attempts to connect to this service with any other
    // version then the service will reply with a status message of UnsupportedVersion and close the client.
    version uint32

    getArchive ArchiveProvider
    commands   chan command
    workers    WorkerPool
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
        if msg.Version != s.version {
            _ = source.SendNow(status.UnsupportedVersion)
            return
        }

        if len(s.sessions) >= s.capacity {
            _ = source.SendNow(status.Full)
            return
        }

        session := s.newSession(source, s.workers)
        s.sessions[source.Id()] = session

        session.OnClosed(func(cli *client.Client) { s.execute(unregisterSession{cli: cli}) })
        session.Process()

        session.Info("Registered client to file service")

        _ = source.SendNow(status.Okay)
    case *PassiveRequest:
        session, exists := s.sessions[source.Id()]
        if !exists {
            source.Fatal(errors.New("fileservice: received request from client that does not have active session"))
            return
        }
        session.EnqueuePassive(msg.Request)
    case *PriorityRequest:
        session, exists := s.sessions[source.Id()]
        if !exists {
            source.Fatal(errors.New("fileservice: received request from client that does not have active session"))
            return
        }
        session.EnqueuePriority(msg.Request)
    }
}

type unregisterSession struct {
    cli *client.Client
}

func (cmd unregisterSession) execute(service *Service) {
    delete(service.sessions, cmd.cli.Id())
    cmd.cli.Info("Unregistered file session")
}
