package account

import "time"

const hoursInADay = 24

// LastLogin is the last time an Account was logged into by an arbitrary user.
type LastLogin time.Time

// DaysSince returns the amount of days that have passed since the value
// of this LastLogin and the value of the given time.Time.
func (lastLogin LastLogin) DaysSince(other time.Time) int {
	return int(time.Time(lastLogin).Sub(other).Hours()) / hoursInADay
}
