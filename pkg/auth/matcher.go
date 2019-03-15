package auth

import (
	"fmt"
	"github.com/sprinkle-it/donut/pkg/account"
	"golang.org/x/crypto/bcrypt"
)

// PasswordMatcher matches a plain Password with its hashed variant to find
// equality in the two inputs. Returns nil on success and an error on failure.
type PasswordMatcher func(plain account.Password, hash account.Password) error

// MatchPasswordsBasic is a basic PasswordMatcher that uses a simple equality
// check to compare the two Password's.
func MatchPasswordsBasic(plain account.Password, hash account.Password) error {
	if plain == hash {
		return nil
	} else {
		return fmt.Errorf("password mismatch")
	}
}

// MatchPasswordsWithBCrypt uses the BCrypt algorithm to match two Password's
// against each other to find equality. Returns nil on success and an error
// on failure.
func MatchPasswordsWithBCrypt(plain account.Password, hash account.Password) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
