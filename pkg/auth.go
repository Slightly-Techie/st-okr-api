package auth

import (
	"fmt"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	googleClientID := config.ENV.GoogleClientID
	googleClientSecret := config.ENV.GoogleClientSecret
	googleCallbackURL := fmt.Sprintf("http://localhost:%s/api/v1/auth/google/callback", config.ENV.ServerPort)

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, googleCallbackURL, "email", "profile"),
	)
}