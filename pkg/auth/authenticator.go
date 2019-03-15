package auth

import "github.com/sprinkle-it/donut/pkg/account"

// AccountSupplier supplies a user Account that is registered
// by the given Email address.
type AccountSupplier func(email account.Email) (*account.Account, error)

// Authenticator authenticates users to see if they are truly
// who they claim to be.
type Authenticator struct {
	supplyAccount  AccountSupplier
	matchPasswords PasswordMatcher
}

func NewAuthenticator(supplier AccountSupplier, passwordMatcher PasswordMatcher) Authenticator {
	return Authenticator{
		supplyAccount:  supplier,
		matchPasswords: passwordMatcher,
	}
}

// Authenticate attempts to authenticate a user using the given email and
// password credentials.
func (auth Authenticator) Authenticate(email account.Email, password account.Password) (Result, error) {
	firstFactorResult, err := auth.doFirstFactor(email, password)
	if err != nil {
		return nil, err
	}

	// only way to type check whilst avoiding reflection unfortunately
	switch result := firstFactorResult.(type) {
	case FirstFactorSuccess:
		return result, nil // TODO second factor

	default:
		return result, nil
	}
}

// doFirstFactor performs the first factor of authentication by looking up the
// associated Account and running a password match against it to ensure a correct
// password input from the user's side. It returns an authentication result which
// might indicate success or failure, or it may return an error which is an
// indication of something very wrong.
func (auth Authenticator) doFirstFactor(email account.Email, password account.Password) (Result, error) {
	accountFetch, err := auth.supplyAccount(email)
	if err != nil {
		return nil, err
	}

	if accountFetch == nil {
		return couldNotFindAccount, nil
	}

	err = auth.matchPasswords(password, accountFetch.Password)
	if err != nil {
		return passwordMismatch, nil
	}

	return FirstFactorSuccess{Account: *accountFetch}, nil
}
