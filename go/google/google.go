package google

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"
)

type Config struct {
	AppName      string
	ClientID     string
	ClientSecret string
}

type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type IDToken struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	AtHash        string `json:"at_hash"`
	HD            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Nonce         string `json:"nonce"`
}

const callbackURL = "http://localhost:3000/callback"

func (c *Config) VerifyIDToken(ctx context.Context, code string) (*IDToken, error) {
	t, err := c.Token(ctx, code)
	if err != nil {
		return nil, err
	}

	return c.VerifyToken(ctx, t.IDToken)
}

func (c *Config) VerifyToken(ctx context.Context, idToken string) (*IDToken, error) {
	v := url.Values{}
	v.Set("id_token", idToken)

	url := "https://oauth3.googleapis.com/tokeninfo?" + v.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var t *IDToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	claims := []string{"https://accounts.google.com", "accounts.google.com"}
	if ok := slices.Contains[[]string, string](claims, t.Iss); !ok {
		return nil, fmt.Errorf("invalid iss: %s", t.Iss)
	}
	if c.ClientID != t.Aud {
		return nil, fmt.Errorf("invalid aud: %s", t.Aud)
	}
	exp, err := strconv.ParseInt(t.Exp, 10, 64)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	if exp < now {
		return nil, fmt.Errorf("expired")
	}

	return t, nil
}

func (c *Config) Token(ctx context.Context, code string) (*TokenInfo, error) {
	v := url.Values{}
	v.Set("code", code)
	v.Set("client_id", c.ClientID)
	v.Set("client_secret", c.ClientSecret)
	v.Set("redirect_uri", callbackURL)
	v.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var t *TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Config) Auth(ctx context.Context) (string, error) {
	const authURL = "https://accounts.google.com/o/oauth2/v2/auth"

	state := generateRandomString(16)
	nonce, err := generateNonce(16)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Set("client_id", c.ClientID)
	v.Set("scope", "openid")
	v.Set("redirect_uri", callbackURL)
	v.Set("state", state)
	v.Set("nonce", nonce)
	v.Set("response_type", "code")

	url := authURL + "?" + v.Encode()

	fmt.Println("url: ", url)

	return url, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func generateNonce(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := crand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
