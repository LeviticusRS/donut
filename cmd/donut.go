package main

import (
    "github.com/sprinkle-it/donut/pkg/asset"
    "github.com/sprinkle-it/donut/pkg/client"
    "github.com/sprinkle-it/donut/pkg/file"
    "github.com/sprinkle-it/donut/pkg/game"
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
        log.Fatal("Failed to open cache: ", err)
    }

    storage, err := asset.NewStorage(cache)
    if err != nil {
        log.Fatal("Failed to create storage: ", err)
    }

    fileService, err := file.New(file.Config{
        LoggerConfig:     loggerConfig,
        Capacity:         1000,
        SupportedVersion: 177,
        Provider:         storage.Get,
        Workers:          2,
        SessionConfig: file.SessionConfig{
            PriorityRequestCapacity: 200,
            PassiveRequestCapacity:  200,
        },
    })

    if err != nil {
        log.Fatal("Failed to create file service")
    }

    fileService.Process()

    gameService, err := game.New(game.Config{

    })

    if err != nil {
        log.Fatal("Failed to create game service")
    }

    gameService.Process()

    srv, err := server.New(server.Config{
        LoggerConfig:   loggerConfig,
        ClientCapacity: 2000,
        ClientConfig:   client.NewDefaultConfig(),
        Receivers: []client.MailReceiver{
            {
                Handler: fileService.HandleMail,
                Accept: []message.Config{
                    file.PassiveRequestConfig,
                    file.PriorityRequestConfig,
                    file.OnlineStatusUpdateConfig,
                    file.OfflineStatusUpdateConfig,
                    file.HandshakeConfig,
                },
            },
            {
                Handler: gameService.HandleMail,
                Accept: []message.Config{
                    game.HandshakeConfig,
                    game.AuthenticateConfig,
                    game.WindowUpdateConfig,
                    game.HeartbeatConfig,
                    game.SceneRebuiltConfig,
                    game.FocusChangedConfig,
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
