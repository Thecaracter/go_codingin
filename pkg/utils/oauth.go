package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string
}

type OAuthUserInfo struct {
	Email      string
	Name       string
	AvatarURL  string
	Provider   string
	ProviderID string
}

func GetGoogleOAuthConfig(config *OAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

func GetGithubOAuthConfig(config *OAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		RedirectURL:  config.GithubRedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func GetGoogleUserInfo(ctx context.Context, code string, config *oauth2.Config) (*OAuthUserInfo, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		AvatarURL:  userInfo.Picture,
		Provider:   "google",
		ProviderID: userInfo.ID,
	}, nil
}

func GetGithubUserInfo(ctx context.Context, code string, config *oauth2.Config) (*OAuthUserInfo, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	// If email is not public, fetch it separately
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailData, err := io.ReadAll(emailResp.Body)
			if err == nil {
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				if json.Unmarshal(emailData, &emails) == nil {
					for _, e := range emails {
						if e.Primary {
							userInfo.Email = e.Email
							break
						}
					}
				}
			}
		}
	}

	if userInfo.Email == "" {
		return nil, errors.New("unable to get email from GitHub")
	}

	name := userInfo.Name
	if name == "" {
		name = userInfo.Login
	}

	return &OAuthUserInfo{
		Email:      userInfo.Email,
		Name:       name,
		AvatarURL:  userInfo.AvatarURL,
		Provider:   "github",
		ProviderID: string(rune(userInfo.ID)),
	}, nil
}
