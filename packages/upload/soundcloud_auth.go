package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const soundcloudAuthorizeURL = "https://secure.soundcloud.com/authorize"
const soundcloudTokenURL = "https://secure.soundcloud.com/oauth/token"

type SoundcloudToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
}

type soundcloudTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type soundcloudPKCEState struct {
	CodeVerifier string `json:"code_verifier"`
	State        string `json:"state"`
}

func soundcloudClientID() string {
	return os.Getenv("SOUNDCLOUD_CLIENT_ID")
}

func soundcloudClientSecret() string {
	return os.Getenv("SOUNDCLOUD_CLIENT_SECRET")
}

func soundcloudRedirectURI() string {
	return os.Getenv("SOUNDCLOUD_REDIRECT_URI")
}

func validateSoundcloudConfig() error {
	switch {
	case soundcloudClientID() == "":
		return fmt.Errorf("SOUNDCLOUD_CLIENT_ID is not set")
	case soundcloudClientSecret() == "":
		return fmt.Errorf("SOUNDCLOUD_CLIENT_SECRET is not set")
	case soundcloudRedirectURI() == "":
		return fmt.Errorf("SOUNDCLOUD_REDIRECT_URI is not set")
	default:
		return nil
	}
}

func soundcloudPKCEPath(tokenPath string) string {
	return fmt.Sprintf("%s.pkce", tokenPath)
}

func randomURLSafeString(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func startSoundcloudAuth(tokenPath string) (string, error) {
	if err := validateSoundcloudConfig(); err != nil {
		return "", err
	}

	codeVerifier, err := randomURLSafeString(32)
	if err != nil {
		return "", err
	}
	state, err := randomURLSafeString(24)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])

	if err := os.MkdirAll(filepath.Dir(tokenPath), 0o755); err != nil {
		return "", err
	}

	pkceState := soundcloudPKCEState{
		CodeVerifier: codeVerifier,
		State:        state,
	}
	pkceJSON, err := json.MarshalIndent(pkceState, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(soundcloudPKCEPath(tokenPath), pkceJSON, 0o600); err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("client_id", soundcloudClientID())
	values.Set("redirect_uri", soundcloudRedirectURI())
	values.Set("response_type", "code")
	values.Set("code_challenge", codeChallenge)
	values.Set("code_challenge_method", "S256")
	values.Set("state", state)

	return fmt.Sprintf("%s?%s", soundcloudAuthorizeURL, values.Encode()), nil
}

func finishSoundcloudAuth(tokenPath string, authCode string) error {
	if err := validateSoundcloudConfig(); err != nil {
		return err
	}

	pkceJSON, err := os.ReadFile(soundcloudPKCEPath(tokenPath))
	if err != nil {
		return fmt.Errorf("failed to read PKCE state: %w", err)
	}

	var pkceState soundcloudPKCEState
	if err := json.Unmarshal(pkceJSON, &pkceState); err != nil {
		return fmt.Errorf("failed to parse PKCE state: %w", err)
	}

	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", soundcloudClientID())
	values.Set("client_secret", soundcloudClientSecret())
	values.Set("redirect_uri", soundcloudRedirectURI())
	values.Set("code_verifier", pkceState.CodeVerifier)
	values.Set("code", authCode)

	token, err := performSoundcloudTokenRequest(values)
	if err != nil {
		return err
	}
	if err := saveSoundcloudToken(tokenPath, token); err != nil {
		return err
	}
	if err := os.Remove(soundcloudPKCEPath(tokenPath)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func refreshSoundcloudToken(tokenPath string) (SoundcloudToken, error) {
	if err := validateSoundcloudConfig(); err != nil {
		return SoundcloudToken{}, err
	}

	token, err := loadSoundcloudToken(tokenPath)
	if err != nil {
		return SoundcloudToken{}, err
	}

	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("client_id", soundcloudClientID())
	values.Set("client_secret", soundcloudClientSecret())
	values.Set("refresh_token", token.RefreshToken)

	refreshedToken, err := performSoundcloudTokenRequest(values)
	if err != nil {
		return SoundcloudToken{}, err
	}
	if err := saveSoundcloudToken(tokenPath, refreshedToken); err != nil {
		return SoundcloudToken{}, err
	}

	return refreshedToken, nil
}

func loadSoundcloudToken(tokenPath string) (SoundcloudToken, error) {
	tokenJSON, err := os.ReadFile(tokenPath)
	if err != nil {
		return SoundcloudToken{}, fmt.Errorf("failed to read SoundCloud token file: %w", err)
	}

	var token SoundcloudToken
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return SoundcloudToken{}, fmt.Errorf("failed to parse SoundCloud token file: %w", err)
	}
	if token.RefreshToken == "" {
		return SoundcloudToken{}, fmt.Errorf("SoundCloud token file is missing refresh_token")
	}

	return token, nil
}

func saveSoundcloudToken(tokenPath string, token SoundcloudToken) error {
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0o755); err != nil {
		return err
	}

	tokenJSON, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath, tokenJSON, 0o600)
}

func performSoundcloudTokenRequest(values url.Values) (SoundcloudToken, error) {
	req, err := http.NewRequest(http.MethodPost, soundcloudTokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return SoundcloudToken{}, err
	}
	req.Header.Set("accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SoundcloudToken{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SoundcloudToken{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return SoundcloudToken{}, fmt.Errorf("SoundCloud token request failed: %s: %s", resp.Status, string(responseBody))
	}

	var tokenResponse soundcloudTokenResponse
	if err := json.Unmarshal(responseBody, &tokenResponse); err != nil {
		return SoundcloudToken{}, fmt.Errorf("failed to parse SoundCloud token response: %w", err)
	}

	return SoundcloudToken{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Scope:        tokenResponse.Scope,
		TokenType:    tokenResponse.TokenType,
	}, nil
}
