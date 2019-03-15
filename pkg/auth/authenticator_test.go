package auth

import (
	"errors"
	"github.com/sprinkle-it/donut/pkg/account"
	"reflect"
	"testing"
)

var dummyAccountSupplier = func(email account.Email) (*account.Account, error) {
	return &account.Account{}, nil
}

var nilAccountSupplier = func(email account.Email) (*account.Account, error) {
	return nil, nil
}

var positivePasswordMatcher = func(plain account.Password, hash account.Password) error {
	return nil
}

var negativePasswordMatcher = func(plain account.Password, hash account.Password) error {
	return errors.New("")
}

func TestAuthenticator_Authenticate(t *testing.T) {
	email := account.Email("sino@gmail.com")
	password := account.Password("sino")

	authenticator := NewAuthenticator(nilAccountSupplier, negativePasswordMatcher)
	result, _ := authenticator.Authenticate(email, password)
	if result != couldNotFindAccount {
		t.Errorf("expected result to be of type CouldNotFindAccount but was %v instead", reflect.TypeOf(result))
	}
}

func TestAuthenticator_Authenticate1(t *testing.T) {
	email := account.Email("sino@gmail.com")
	password := account.Password("sino")

	authenticator := NewAuthenticator(dummyAccountSupplier, negativePasswordMatcher)
	result, _ := authenticator.Authenticate(email, password)
	if result != passwordMismatch {
		t.Errorf("expected result to be of type PasswordMismatch but was %v instead", reflect.TypeOf(result))
	}
}

func TestAuthenticator_Authenticate2(t *testing.T) {
	email := account.Email("sino@gmail.com")
	password := account.Password("sino")

	authenticator := NewAuthenticator(dummyAccountSupplier, positivePasswordMatcher)
	result, _ := authenticator.Authenticate(email, password)

	isSuccess := reflect.TypeOf(result) == reflect.TypeOf(FirstFactorSuccess{})
	if !isSuccess {
		t.Errorf("expected result to be of type FirstFactorSuccess but was %v instead", reflect.TypeOf(result))
	}
}
