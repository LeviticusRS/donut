package game

type Config struct {
    Capacity    int
    WorldConfig WorldConfig
}

func (c Config) Build() (*Service, error) {
    return &Service{
        capacity: c.Capacity,
        commands: make(chan command),
        world:    c.WorldConfig.Build(),
    }, nil
}