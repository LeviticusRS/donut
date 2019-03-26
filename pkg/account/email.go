package account

// Email is the e-mail address that is associated with an Account.
type Email string

// IsValid returns whether the value of this Email is
// valid according to the pre-defined domain rules.
func (email Email) IsValid() bool {
	return true // TODO validation
}
