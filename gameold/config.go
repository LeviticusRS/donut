package gameold

type Config struct {
    Capacity    int
    Authenticator Authenticator
    WorldConfig WorldConfig
}

func (c Config) Build() (*Service, error) {
    return &Service{
        capacity: c.Capacity,
        authenticator: c.Authenticator,
        commands: make(chan command),
        world:    c.WorldConfig.Build(),
    }, nil
}
