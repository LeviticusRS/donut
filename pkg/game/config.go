package game

type Config struct {
    Capacity int
}

func (c Config) Build() (*Service, error) {
    return &Service{
        capacity: c.Capacity,
        commands: make(chan command),
    }, nil
}