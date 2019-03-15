package account

import "testing"

func TestPassword_IsValid(t *testing.T) {
	password := Password("hello")
	if password.IsValid() {
		t.Error("expected password to be invalid as it has less than six characters")
	}
}

func TestPassword_IsValid2(t *testing.T) {
	password := Password("hello123")
	if !password.IsValid() {
		t.Error("expected password to be valid as it has more than six characters")
	}
}
