package auth

import (
	"github.com/sprinkle-it/donut/pkg/account"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestMatchPasswordsWithBCrypt(t *testing.T) {
	plain := account.Password("praisesino")

	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	hash := account.Password(hashBytes)

	passwordsMatch := MatchPasswordsWithBCrypt(plain, hash) == nil
	if !passwordsMatch {
		t.Error("expected passwords to match")
	}
}

func TestMatchPasswordsWithBCrypt2(t *testing.T) {
	plain := account.Password("praisesino")
	plain2 := account.Password("dontpraisesini")

	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(plain2), bcrypt.DefaultCost)
	hash := account.Password(hashBytes)

	passwordsMatch := MatchPasswordsWithBCrypt(plain, hash) == nil
	if passwordsMatch {
		t.Error("expected passwords to mismatch")
	}
}
