package helpers

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Access token expires in 24 hours
	})

	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func GenerateRefreshToken() (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func FormatTextToSlug(text string) string {

	text = strings.ToLower(text)

	// Trim whitespace from both ends
	text = strings.TrimSpace(text)

	// Remove all non-alphanumeric characters except spaces
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	text = reg.ReplaceAllString(text, "")

	// Replace spaces and hyphens with single hyphens
	text = strings.ReplaceAll(text, " ", "-")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	text = reg.ReplaceAllString(text, "-")

	return text
}
