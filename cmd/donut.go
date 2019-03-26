package main

import (
	"github.com/sprinkle-it/donut/pkg/account"
	"github.com/sprinkle-it/donut/pkg/asset"
    "github.com/sprinkle-it/donut/pkg/file"
    "github.com/sprinkle-it/donut/pkg/game"
    "github.com/sprinkle-it/donut/pkg/message"
    "github.com/sprinkle-it/donut/server"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "log"
    "net/http"
    _ "net/http/pprof"
    "time"
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

	accountRepository := account.NewDummyRepository()
	game.NewAuthenticator( // TODO
		game.SupplyAccountFromRepository(accountRepository),
		account.MatchPasswordsBasic,
	)

	gameService, err := game.New(game.Config{
		WorldConfig: game.WorldConfig{
			PlayerCapacity: 2048,
			Delta:          time.Millisecond * 600,
		},
	})

	if err != nil {
		log.Fatal("Failed to create game service")
	}

	gameService.Process()

	srv, err := server.New(server.Config{
		LoggerConfig:   loggerConfig,
		ClientCapacity: 2000,
		ClientConfig:   server.NewDefaultClientConfig(),
		Receivers: []server.MailReceiver{
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
					game.NewLoginConfig,
					game.WindowUpdateConfig,
					game.HeartbeatConfig,
					game.ClientPerformanceMeasuredConfig,
					game.MouseActivityRecordedConfig,
					game.MouseClickedConfig,
					game.MinimapWalkConfig,
					game.WalkHereConfig,
					game.ExamineObjectConfig,
					game.ButtonPressedConfig,
					game.SceneRebuiltConfig,
					game.FocusChangedConfig,
					game.KeyTypedConfig,
					game.CameraRotatedConfig,
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
