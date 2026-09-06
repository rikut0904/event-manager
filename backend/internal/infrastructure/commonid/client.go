package commonid

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidCallback = errors.New("Common IDのコールバックが不正です")
	ErrStateMismatch   = errors.New("Common IDのstateが一致しません")
	ErrInvalidSession  = errors.New("Common IDのセッションが無効です")
)

type Config struct {
	Origin            string
	APIOrigin         string
	ClientID          string
	RedirectURI       string
	LogoutRedirectURI string
	APIKey            string
}

type Pending struct {
	State       string
	Verifier    string
	ClientID    string
	RedirectURI string
	ExpiresAt   time.Time
}

type User struct {
	CommonUserID string `json:"common_user_id"`
	Email        string `json:"email"`
}

type Client struct {
	cfg Config
}

func New(cfg Config) (*Client, error) {
	if cfg.Origin == "" || cfg.APIOrigin == "" || cfg.ClientID == "" || cfg.RedirectURI == "" || cfg.LogoutRedirectURI == "" || cfg.APIKey == "" {
		return nil, errors.New("Common IDのOrigin、API Origin、client_id、リダイレクトURI、APIキーが必要です")
	}
	return &Client{cfg: cfg}, nil
}

func (c *Client) Begin(intent string) (string, Pending, error) {
	if intent != "login" && intent != "signup" {
		return "", Pending{}, errors.New("認証 intent が不正です")
	}
	state, err := randomString()
	if err != nil {
		return "", Pending{}, err
	}
	verifier, err := randomString()
	if err != nil {
		return "", Pending{}, err
	}
	pending := Pending{State: state, Verifier: verifier, ClientID: c.cfg.ClientID, RedirectURI: c.cfg.RedirectURI, ExpiresAt: time.Now().UTC().Add(10 * time.Minute)}
	query := url.Values{
		"client_id": {c.cfg.ClientID}, "redirect_uri": {c.cfg.RedirectURI}, "response_type": {"code"},
		"scope": {"openid profile email"}, "state": {state}, "code_challenge": {challenge(verifier)},
		"code_challenge_method": {"S256"}, "intent": {intent},
	}
	return strings.TrimRight(c.cfg.Origin, "/") + "/auth?" + query.Encode(), pending, nil
}

func (c *Client) BeginLogout() (string, string, error) {
	state, err := randomString()
	if err != nil {
		return "", "", err
	}
	query := url.Values{"client_id": {c.cfg.ClientID}, "post_logout_redirect_uri": {c.cfg.LogoutRedirectURI}, "state": {state}}
	return strings.TrimRight(c.cfg.Origin, "/") + "/logout?" + query.Encode(), state, nil
}

func (c *Client) Exchange(ctx context.Context, callback url.Values, pending Pending) (User, error) {
	if pending.State == "" || pending.Verifier == "" || pending.ClientID != c.cfg.ClientID || pending.RedirectURI != c.cfg.RedirectURI || time.Now().UTC().After(pending.ExpiresAt) {
		return User{}, ErrInvalidCallback
	}
	if callback.Get("state") != pending.State {
		return User{}, ErrStateMismatch
	}
	if callback.Get("code") == "" || callback.Get("error") != "" {
		return User{}, ErrInvalidCallback
	}
	form := url.Values{"grant_type": {"authorization_code"}, "code": {callback.Get("code")}, "client_id": {pending.ClientID}, "redirect_uri": {pending.RedirectURI}, "code_verifier": {pending.Verifier}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.APIOrigin, "/")+"/v1/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-Key", c.cfg.APIKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		io.Copy(io.Discard, res.Body)
		return User{}, ErrInvalidCallback
	}
	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return User{}, err
	}
	if user.CommonUserID == "" {
		return User{}, ErrInvalidCallback
	}
	return user, nil
}

func (c *Client) CheckSession(ctx context.Context, token string) (User, error) {
	if token == "" {
		return User{}, ErrInvalidSession
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.cfg.APIOrigin, "/")+"/v1/oauth/session", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("X-API-Key", c.cfg.APIKey)
	req.Header.Set("X-Client-ID", c.cfg.ClientID)
	req.AddCookie(&http.Cookie{Name: "common_id_session", Value: token})
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return User{}, ErrInvalidSession
	}
	var user struct {
		CommonUserID string `json:"common_user_id"`
		Email        string `json:"email"`
		Status       string `json:"status"`
	}
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return User{}, err
	}
	if user.CommonUserID == "" || user.Status != "active" {
		return User{}, ErrInvalidSession
	}
	return User{CommonUserID: user.CommonUserID, Email: user.Email}, nil
}

func randomString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("乱数生成に失敗しました: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func challenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
