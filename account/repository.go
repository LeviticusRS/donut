package account

// Repository is the storage of user Account's.
type Repository interface {
	Get(email Email) (*Account, error)
	Put(email Email, account Account) error
}

// DummyRepository is an implementation of Repository
// that supplies dummy Account instances.
type DummyRepository struct{}

func NewDummyRepository() *DummyRepository {
	return &DummyRepository{}
}

func (repository *DummyRepository) Get(email Email) (*Account, error) {
	return &Account{
		Email:       email,
		Password:    Password("hello123"),
		DisplayName: DisplayName("sino"),
	}, nil
}

func (repository *DummyRepository) Put(email Email, account Account) error {
	return nil
}
