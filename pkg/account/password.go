package account

import "golang.org/x/crypto/bcrypt"

const (
	// TODO figure out what an appropriate cost is. the default is 10. what about generating a cost per password??
	passwordHashCost = 12

	minimalPasswordLength = 6
)

// Password is a variable length partial secret that only the user of
// an account may know to gain access to their account.
type Password string

// PasswordHasher accepts a plain Password and returns
// its hashed variant.
type PasswordHasher func(plain Password) (Password, error)

// HashPasswordWithBCrypt hashes the given plain Password
// using the BCrypt algorithm.
func HashPasswordWithBCrypt(plain Password) (Password, error) {
	hashValue, err := bcrypt.GenerateFromPassword([]byte(plain), passwordHashCost)
	if err != nil {
		return Password(""), err
	}

	return Password(hashValue), nil
}

// IsValid returns whether the value of this Password is valid
// according to the pre-defined domain rules.
func (password Password) IsValid() bool {
	return len(password) >= minimalPasswordLength // TODO further validation
}
