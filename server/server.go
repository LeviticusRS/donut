package server

import (
    "fmt"
    "go.uber.org/zap"
    "net"
)

type Config struct {
    LoggerConfig   zap.Config
    ClientCapacity int
    ClientConfig   ClientConfig
    Receivers      []MailReceiver
}

func (cfg Config) Build() (*Server, error) {
    logger, err := cfg.LoggerConfig.Build()
    if err != nil {
        return nil, err
    }

    router, err := NewMailRouter(cfg.Receivers)
    if err != nil {
        return nil, err
    }

    return &Server{
        logger:         logger,
        clientCapacity: cfg.ClientCapacity,
        clientFactory:  cfg.ClientConfig.Build,
        clients:        make(map[uint64]*Client, cfg.ClientCapacity),
        router:         router,
        commands:       make(chan command),
    }, nil
}

type Server struct {
    logger *zap.Logger

    clientCapacity int
    clientFactory  Factory
    clients        map[uint64]*Client
    router         MailRouter

    commands chan command
}

func New(config Config) (*Server, error) {
    return config.Build()
}

func (s *Server) execute(cmd command) { s.commands <- cmd }

func (s *Server) process() {
    go func() {
        for cmd := range s.commands {
            cmd.Execute(s)
        }
    }()
}

func (s *Server) Listen(port int) error {
    listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
    if err != nil {
        return err
    }

    s.logger.Info("Listening", zap.Int("port", port))

    s.process()

    for {
        conn, err := listener.Accept()
        if err != nil {
            return err
        }
        s.execute(acceptConnection{connection: conn,})
    }
}

// Server commands are calls that are used to assure that multiple processes needing to update the server state
// can be done so synchronously.
type command interface {
    Execute(srv *Server)
}

// Accepts the wrapped connection and creates a new client. If the server is at capacity the connection will be closed.
type acceptConnection struct {
    connection net.Conn
}

func (cmd acceptConnection) Execute(server *Server) {
    if len(server.clients) >= server.clientCapacity {
        server.logger.Info("Failed to accept, server is at capacity",
            zap.Stringer("address", cmd.connection.RemoteAddr()),
        )
        _ = cmd.connection.Close()
        return
    }

    cli := server.clientFactory(cmd.connection, server.logger, server.router)
    server.clients[cli.Id()] = cli

    // Register a callback that will unregister the client from the server when it is closed.
    cli.OnClosed(func(cli *Client) {
        server.execute(unregisterClient{client: cli,})
    })

    server.logger.Info("Registered client",
        zap.Uint64("id", cli.Id()),
        zap.Stringer("address", cli.RemoteAddress()),
    )

    cli.Process()
}

// Unregisters a client from the server freeing the reference in the server.
type unregisterClient struct {
    client *Client
}

func (cmd unregisterClient) Execute(server *Server) {
    delete(server.clients, cmd.client.Id())

    server.logger.Info("Unregistered client",
        zap.Uint64("id", cmd.client.Id()),
        zap.Stringer("address", cmd.client.RemoteAddress()),
    )
}
