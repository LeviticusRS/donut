package auth

import "github.com/sprinkle-it/donut/pkg/account"

var (
	passwordMismatch    = PasswordMismatch{}
	couldNotFindAccount = CouldNotFindAccount{}
)

// FirstFactorSuccess is an indication of the first factor procedure
// having successfully been authenticated.
type FirstFactorSuccess struct {
	Account account.Account
}

// PasswordMismatch is an authentication Result that indicates
// two given Password's did not match, meaning the user has
// entered an invalid password.
type PasswordMismatch struct{}

// CouldNotFindAccount is an authentication Result that indicates
// that a user does not exist in the database.
type CouldNotFindAccount struct{}

// Result is the result from attempting to authenticate a user.
type Result interface{}
