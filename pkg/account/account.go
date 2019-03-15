package account

// Account represents the account that a user has registered.
type Account struct {
	Email               Email
	Password            Password
	DisplayName         DisplayName
	PreviousDisplayName *DisplayName
	LastLogin           *LastLogin
}
