package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleAuthService struct {
	config *oauth2.Config
}

func NewGoogleAuthService() *GoogleAuthService {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleAuthService{
		config: config,
	}
}

// GetAuthURL genera la URL para iniciar el flujo de OAuth2
func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

// Exchange intercambia el código de autorización por un token
func (s *GoogleAuthService) Exchange(code string) (*oauth2.Token, error) {
	return s.config.Exchange(oauth2.NoContext, code)
}

// GetUserInfo obtiene la información del usuario de Google
func (s *GoogleAuthService) GetUserInfo(token *oauth2.Token) (*GoogleUser, error) {
	client := s.config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error al obtener información del usuario")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GoogleUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}