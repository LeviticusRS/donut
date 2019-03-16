package game

type Config struct {
	Capacity int

	Authenticator Authenticator
}

func (c Config) Build() (*Service, error) {
	return &Service{
		capacity:      c.Capacity,
		authenticator: c.Authenticator,
		commands:      make(chan command),
	}, nil
}
