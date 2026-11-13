package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

const (
	PrefGoogleUserName  = "google_user_name"
	PrefGoogleUserEmail = "google_user_email"
	PrefGoogleAuthToken = "google_auth_token"
)

// GoogleUser stores the authenticated user's profile information
type GoogleUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	authServerMu sync.Mutex
	authServer   *http.Server
)

// GetSignedInUser returns the currently stored Google user profile, if any
func GetSignedInUser(prefs fyne.Preferences) (*GoogleUser, bool) {
	if prefs == nil {
		return nil, false
	}
	email := strings.TrimSpace(prefs.String(PrefGoogleUserEmail))
	if email == "" {
		return nil, false
	}
	name := prefs.StringWithFallback(PrefGoogleUserName, "Google User")
	return &GoogleUser{
		Name:  name,
		Email: email,
	}, true
}

// SaveSignedInUser stores user profile information in preferences
func SaveSignedInUser(prefs fyne.Preferences, name, email, token string) {
	if prefs == nil {
		return
	}
	prefs.SetString(PrefGoogleUserName, name)
	prefs.SetString(PrefGoogleUserEmail, email)
	if token != "" {
		prefs.SetString(PrefGoogleAuthToken, token)
	}
}

// SignOut clears Google authentication credentials
func SignOut(prefs fyne.Preferences) {
	if prefs == nil {
		return
	}
	prefs.SetString(PrefGoogleUserName, "")
	prefs.SetString(PrefGoogleUserEmail, "")
	prefs.SetString(PrefGoogleAuthToken, "")
}

// StartOAuthLoopbackListener creates a local callback listener on 127.0.0.1:8085 for OAuth redirects
func StartOAuthLoopbackListener(onCodeReceived func(code string, err error)) error {
	authServerMu.Lock()
	defer authServerMu.Unlock()

	if authServer != nil {
		_ = authServer.Close()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		errStr := r.URL.Query().Get("error")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if errStr != "" {
			fmt.Fprintf(w, "<html><body><h3>Google Sign-In Failed: %s</h3><p>You can close this window.</p></body></html>", errStr)
			if onCodeReceived != nil {
				onCodeReceived("", fmt.Errorf("oauth error: %s", errStr))
			}
			return
		}

		fmt.Fprintf(w, "<html><body><h3>Google Sign-In Successful!</h3><p>You can close this tab and return to Mendu: IdeaMapper.</p></body></html>")
		if onCodeReceived != nil {
			onCodeReceived(code, nil)
		}
	})

	authServer = &http.Server{
		Addr:         "127.0.0.1:8085",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		_ = authServer.ListenAndServe()
	}()

	return nil
}

