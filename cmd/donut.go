package main

import (
    "github.com/sprinkle-it/coffee"
    "github.com/sprinkle-it/donut/file"
    "github.com/sprinkle-it/donut/game"
    "github.com/sprinkle-it/donut/server"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
)

func main() {
    loggerConfig := zap.NewDevelopmentConfig()
    loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    loggerConfig.DisableCaller = true

    cache, err := coffee.OpenCache("cache")
    if err != nil {
        log.Fatal("Failed to open cache: ", err)
    }

    storage, err := coffee.NewStorage(cache)
    if err != nil {
        log.Fatal("Failed to create storage: ", err)
    }

    fileService, err := file.New(file.Config{
        LoggerConfig:     loggerConfig,
        Capacity:         1000,
        SupportedVersion: 177,
        ArchiveProvider:  storage.GetArchive,
        Workers:          2,
        SessionConfig: file.SessionConfig{
            PriorityRequestCapacity: 200,
            PassiveRequestCapacity:  200,
        },
    })

    if err != nil {
        log.Fatal("Failed to create file service: ", err)
    }

    fileService.Process()

    gameService, err := game.New(game.Config{
        LoggerConfig:     loggerConfig,
        SupportedVersion: 177,
    })

    if err != nil {
        log.Fatal("Failed to create game service: ", err)
    }

    gameService.Process()

    srv, err := server.New(server.Config{
        LoggerConfig:   loggerConfig,
        ClientCapacity: 2000,
        ClientConfig:   server.NewDefaultClientConfig(),
        Receivers: []server.MailReceiver{
            fileService.MailReceiver(),
            gameService.MailReceiver(),
        },
    })

    if err != nil {
        log.Fatal("Failed to create server: ", err)
    }

    if err := srv.Listen(43594); err != nil {
        log.Fatal("Failed to listen to server port: ", err)
    }
}
