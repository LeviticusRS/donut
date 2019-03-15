package account

import "testing"

func TestDisplayName_IsValid(t *testing.T) {
	displayName := Password("he")
	if displayName.IsValid() {
		t.Error("expected display name to be invalid as it has less than three characters")
	}
}

func TestDisplayName_IsValid2(t *testing.T) {
	displayName := DisplayName("hello")
	if !displayName.IsValid() {
		t.Error("expected display name to be valid as it has more than three characters")
	}
}
