package main

import (
    "github.com/sprinkle-it/donut/pkg/asset"
    "github.com/sprinkle-it/donut/pkg/client"
    "github.com/sprinkle-it/donut/pkg/fileservice"
    "github.com/sprinkle-it/donut/pkg/message"
    "github.com/sprinkle-it/donut/pkg/server"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    loggerConfig := zap.NewDevelopmentConfig()
    loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    loggerConfig.DisableCaller = true

    cache, err := asset.OpenCache("cache", asset.IndexCount)
    if err != nil {
        log.Fatal("Failed to open cache")
    }

    storage, err := asset.NewStorage(cache)
    if err != nil {
        log.Fatal("Failed to create storage")
    }

    service, err := fileservice.New(fileservice.Config{
        LoggerConfig:     loggerConfig,
        Capacity:         1000,
        SupportedVersion: 177,
        Provider:         storage.Get,
        Workers:          2,
        SessionConfig: fileservice.SessionConfig{
            PriorityRequestCapacity: 200,
            PassiveRequestCapacity:  200,
        },
    })
    service.Process()

    if err != nil {
        log.Fatal("Failed to create file service")
    }

    srv, err := server.New(server.Config{
        LoggerConfig:   loggerConfig,
        ClientCapacity: 2000,
        ClientConfig:   client.NewDefaultConfig(),
        Receivers: []client.MailReceiver{
            {
                Handler: service.HandleMail,
                Accept: []message.Descriptor{
                    fileservice.PassiveRequestDescriptor,
                    fileservice.PriorityRequestDescriptor,
                    fileservice.OnlineStatusUpdateDescriptor,
                    fileservice.OfflineStatusUpdateDescriptor,
                    fileservice.HandshakeDescriptor,
                },
            },
        },
    })

    if err != nil {
        log.Fatal("Failed to create server", err)
    }

    if err := srv.Listen(43594); err != nil {
        log.Fatal("Failed to listen to server port", err)
    }
}
