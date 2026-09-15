package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	CookieName                 = "app_session"
	OAuthPendingCookieName     = "common_id_pending"
	OAuthLogoutStateCookieName = "common_id_logout_state"
	defaultTTL                 = 24 * time.Hour
)

var ErrInvalid = errors.New("アプリセッションが無効です")

type Manager struct {
	secret []byte
	secure bool
	ttl    time.Duration
}

type claims struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type OAuthPending struct {
	State       string `json:"state"`
	Verifier    string `json:"verifier"`
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	BackPath    string `json:"back_path"`
	ExpiresAt   int64  `json:"expires_at"`
}

type OAuthLogoutState struct {
	State     string `json:"state"`
	ExpiresAt int64  `json:"expires_at"`
}

func New(secret string, secure bool) (*Manager, error) {
	secret = strings.TrimSpace(secret)
	if len(secret) < 32 {
		return nil, errors.New("APP_SESSION_SECRETは32文字以上で設定してください")
	}
	return &Manager{secret: []byte(secret), secure: secure, ttl: defaultTTL}, nil
}

func (m *Manager) Issue(userID string) (*http.Cookie, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalid
	}
	payload, err := json.Marshal(claims{UserID: userID, ExpiresAt: time.Now().Add(m.ttl).Unix()})
	if err != nil {
		return nil, err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return &http.Cookie{Name: CookieName, Value: encodedPayload + "." + m.sign(encodedPayload), Path: "/", HttpOnly: true, Secure: m.secure, SameSite: http.SameSiteLaxMode, MaxAge: int(m.ttl.Seconds()), Expires: time.Now().Add(m.ttl)}, nil
}

func (m *Manager) Verify(value string) (string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(parts[1]), []byte(m.sign(parts[0]))) {
		return "", ErrInvalid
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalid
	}
	var data claims
	if json.Unmarshal(payload, &data) != nil || data.UserID == "" || data.ExpiresAt <= time.Now().Unix() {
		return "", ErrInvalid
	}
	return data.UserID, nil
}

func (m *Manager) ClearCookie() *http.Cookie {
	return &http.Cookie{Name: CookieName, Value: "", Path: "/", HttpOnly: true, Secure: m.secure, SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)}
}

func (m *Manager) IssueOAuthPending(data OAuthPending) (*http.Cookie, error) {
	if data.State == "" || data.Verifier == "" || data.ClientID == "" || data.RedirectURI == "" || data.ExpiresAt <= time.Now().Unix() {
		return nil, ErrInvalid
	}
	return m.issueSignedCookie(OAuthPendingCookieName, data, time.Until(time.Unix(data.ExpiresAt, 0)))
}

func (m *Manager) VerifyOAuthPending(value string) (OAuthPending, error) {
	var data OAuthPending
	if err := m.verifySigned(value, &data); err != nil || data.State == "" || data.Verifier == "" || data.ClientID == "" || data.RedirectURI == "" || data.ExpiresAt <= time.Now().Unix() {
		return OAuthPending{}, ErrInvalid
	}
	return data, nil
}

func (m *Manager) ClearOAuthPendingCookie() *http.Cookie {
	return m.clearCookie(OAuthPendingCookieName)
}

func (m *Manager) IssueOAuthLogoutState(state string, expiresAt time.Time) (*http.Cookie, error) {
	if state == "" || !expiresAt.After(time.Now()) {
		return nil, ErrInvalid
	}
	data := OAuthLogoutState{State: state, ExpiresAt: expiresAt.Unix()}
	return m.issueSignedCookie(OAuthLogoutStateCookieName, data, time.Until(expiresAt))
}

func (m *Manager) VerifyOAuthLogoutState(value string) (OAuthLogoutState, error) {
	var data OAuthLogoutState
	if err := m.verifySigned(value, &data); err != nil || data.State == "" || data.ExpiresAt <= time.Now().Unix() {
		return OAuthLogoutState{}, ErrInvalid
	}
	return data, nil
}

func (m *Manager) ClearOAuthLogoutStateCookie() *http.Cookie {
	return m.clearCookie(OAuthLogoutStateCookieName)
}

func (m *Manager) issueSignedCookie(name string, value any, ttl time.Duration) (*http.Cookie, error) {
	if ttl <= 0 {
		return nil, ErrInvalid
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	expires := time.Now().Add(ttl)
	return &http.Cookie{Name: name, Value: encodedPayload + "." + m.sign(encodedPayload), Path: "/", HttpOnly: true, Secure: m.secure, SameSite: http.SameSiteLaxMode, MaxAge: int(ttl.Seconds()), Expires: expires}, nil
}

func (m *Manager) verifySigned(value string, target any) error {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(parts[1]), []byte(m.sign(parts[0]))) {
		return ErrInvalid
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil || json.Unmarshal(payload, target) != nil {
		return ErrInvalid
	}
	return nil
}

func (m *Manager) clearCookie(name string) *http.Cookie {
	return &http.Cookie{Name: name, Value: "", Path: "/", HttpOnly: true, Secure: m.secure, SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)}
}

func (m *Manager) sign(value string) string {
	hash := hmac.New(sha256.New, m.secret)
	_, _ = hash.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
