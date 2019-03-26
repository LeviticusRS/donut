package account

import "time"

const (
    minimumDisplayNameLength = 3
    hoursInADay = 24
)

// Email is the e-mail address that is associated with an Account.
type Email string

// IsValid returns whether the value of this Email is valid according to the pre-defined domain rules.
func (email Email) IsValid() bool {
    return true // TODO validation
}

// DisplayName is the name of a user account that is exposed to others in-game for identification purposes.
type DisplayName string

// IsValid returns whether the value of this DisplayName is valid according to the pre-defined domain rules.
func (name DisplayName) IsValid() bool {
    return len(name) >= minimumDisplayNameLength
}

// LastLogin is the last time an Account was logged into by an arbitrary user.
type LastLogin time.Time

// DaysSince returns the amount of days that have passed since the value of this LastLogin and the value of the
// given time.Time.
func (lastLogin LastLogin) DaysSince(other time.Time) int {
    return int(time.Time(lastLogin).Sub(other).Hours()) / hoursInADay
}

// Account represents the account that a user has registered.
type Account struct {
    Email               Email
    Password            Password
    DisplayName         DisplayName
    PreviousDisplayName *DisplayName
    LastLogin           *LastLogin
}
