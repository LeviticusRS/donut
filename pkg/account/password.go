package account

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

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

// PasswordMatcher matches a plain Password with its hashed variant to find
// equality in the two inputs. Returns nil on success and an error on failure.
type PasswordMatcher func(plain Password, hash Password) error

// MatchPasswordsBasic is a basic PasswordMatcher that uses a simple equality
// check to compare the two Password's.
func MatchPasswordsBasic(plain Password, hash Password) error {
	if plain == hash {
		return nil
	} else {
		return fmt.Errorf("password mismatch")
	}
}

// MatchPasswordsWithBCrypt uses the BCrypt algorithm to match two Password's
// against each other to find equality. Returns nil on success and an error
// on failure.
func MatchPasswordsWithBCrypt(plain Password, hash Password) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

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
