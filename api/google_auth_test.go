package api

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestGoogleAuthPreferences(t *testing.T) {
	app := test.NewApp()
	prefs := app.Preferences()

	// 1. Initial State -> Not signed in
	user, signedIn := GetSignedInUser(prefs)
	if signedIn || user != nil {
		t.Fatalf("expected signedIn to be false initially")
	}

	// 2. Save user
	SaveSignedInUser(prefs, "Aditya B", "aditya@example.com", "fake_token_123")

	user, signedIn = GetSignedInUser(prefs)
	if !signedIn || user == nil {
		t.Fatalf("expected signedIn to be true after SaveSignedInUser")
	}
	if user.Name != "Aditya B" || user.Email != "aditya@example.com" {
		t.Fatalf("unexpected user details: %+v", user)
	}

	// 3. Sign Out
	SignOut(prefs)

	user, signedIn = GetSignedInUser(prefs)
	if signedIn || user != nil {
		t.Fatalf("expected signedIn to be false after SignOut")
	}
}
