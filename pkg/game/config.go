package game

import "github.com/sprinkle-it/donut/pkg/auth"

type Config struct {
	Capacity int

	Authenticator auth.Authenticator
}

func (c Config) Build() (*Service, error) {
	return &Service{
		capacity:      c.Capacity,
		authenticator: c.Authenticator,
		commands:      make(chan command),
	}, nil
}
