package account

const minimalDisplayNameLength = 3

// DisplayName is the name of a user account that is exposed
// to others ingame for identification purposes.
type DisplayName string

// IsValid returns whether the value of this DisplayName is
// valid according to the pre-defined domain rules.
func (name DisplayName) IsValid() bool {
	return len(name) >= minimalDisplayNameLength
}
