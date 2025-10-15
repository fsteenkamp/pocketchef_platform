package auth

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

const COOKIE = "AUTH_TOKEN"

const SessionDuration = time.Hour * 730

func InitGoogle(clientID string, clientSecret string, redirectHost string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/auth/callback/google", redirectHost),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func InitMicrosoft(clientID string, clientSecret string, redirectHost string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/auth/callback/microsoft", redirectHost),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"User.Read",
		},

		// Endpoint: oauth2.Endpoint{
		// 	AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		// 	TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		// },
		Endpoint: microsoft.LiveConnectEndpoint, // would this work? URLs look weird
	}
}
