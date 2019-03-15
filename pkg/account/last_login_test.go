package account

import (
	"testing"
	"time"
)

func TestLastLogin_DaysSince(t *testing.T) {
	today := time.Now()
	lastLogin := LastLogin(today.Add(time.Hour * 48))

	daysDiff := lastLogin.DaysSince(today)
	if daysDiff != 2 {
		t.Errorf("expected the amount of days to equal 2 but was %v instead\n", daysDiff)
	}
}
