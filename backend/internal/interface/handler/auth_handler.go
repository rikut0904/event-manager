package handler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend/internal/infrastructure/commonid"
	"backend/internal/infrastructure/session"
	"backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	commonID    *commonid.Client
	appSession  *session.Manager
	frontendURL string
}

func NewAuthHandler(u usecase.AuthUsecase, commonIDClient *commonid.Client, frontendOrigin string, appSession *session.Manager) *AuthHandler {
	return &AuthHandler{authUsecase: u, commonID: commonIDClient, appSession: appSession, frontendURL: strings.TrimRight(frontendOrigin, "/")}
}

func (h *AuthHandler) Begin(c echo.Context) error {
	if h.commonID == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Common IDが設定されていません")
	}
	intent := c.Param("intent")
	authURL, pending, err := h.commonID.Begin(intent)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pendingCookie, err := h.appSession.IssueOAuthPending(session.OAuthPending{
		State:       pending.State,
		Verifier:    pending.Verifier,
		ClientID:    pending.ClientID,
		RedirectURI: pending.RedirectURI,
		BackPath:    safeBackPath(c.QueryParam("back_path")),
		ExpiresAt:   pending.ExpiresAt.Unix(),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "認証状態を保存できません")
	}
	c.SetCookie(pendingCookie)
	return c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) Callback(c echo.Context) error {
	if h.commonID == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Common IDが設定されていません")
	}
	pendingCookie, err := c.Cookie(session.OAuthPendingCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "認証状態が見つかりません。もう一度お試しください")
	}
	pendingData, err := h.appSession.VerifyOAuthPending(pendingCookie.Value)
	c.SetCookie(h.appSession.ClearOAuthPendingCookie())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "認証状態の有効期限が切れています。もう一度お試しください")
	}
	pending := commonid.Pending{State: pendingData.State, Verifier: pendingData.Verifier, ClientID: pendingData.ClientID, RedirectURI: pendingData.RedirectURI, ExpiresAt: time.Unix(pendingData.ExpiresAt, 0)}
	commonUser, err := h.commonID.Exchange(c.Request().Context(), c.QueryParams(), pending)
	if err != nil {
		log.Printf("Common ID token exchange failed: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Common IDとの認証に失敗しました")
	}
	if _, err := h.authUsecase.SyncCommonUser(c.Request().Context(), commonUser.CommonUserID, commonUser.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザー情報を保存できませんでした")
	}
	appCookie, err := h.appSession.Issue(commonUser.CommonUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "アプリセッションを発行できませんでした")
	}
	c.SetCookie(appCookie)
	backPath := pendingData.BackPath
	if backPath == "" {
		backPath = "/home"
	}
	return c.Redirect(http.StatusFound, h.frontendRedirect(backPath))
}

func (h *AuthHandler) AuthCallback(c echo.Context) error {
	return h.Callback(c)
}

func (h *AuthHandler) BeginLogout(c echo.Context) error {
	c.SetCookie(h.appSession.ClearCookie())
	if h.commonID == nil {
		return c.Redirect(http.StatusFound, h.frontendRedirect("/"))
	}
	logoutURL, state, err := h.commonID.BeginLogout()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ログアウトを開始できません")
	}
	logoutCookie, err := h.appSession.IssueOAuthLogoutState(state, time.Now().UTC().Add(10*time.Minute))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ログアウト状態を保存できません")
	}
	c.SetCookie(logoutCookie)
	return c.Redirect(http.StatusFound, logoutURL)
}

func (h *AuthHandler) LogoutCallback(c echo.Context) error {
	c.SetCookie(h.appSession.ClearCookie())
	logoutCookie, err := c.Cookie(session.OAuthLogoutStateCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ログアウトを確認できません")
	}
	expected, err := h.appSession.VerifyOAuthLogoutState(logoutCookie.Value)
	c.SetCookie(h.appSession.ClearOAuthLogoutStateCookie())
	if err != nil || c.QueryParam("state") != expected.State || (c.QueryParam("logout") != "success" && c.QueryParam("result") != "success") {
		return echo.NewHTTPError(http.StatusBadRequest, "ログアウトを確認できません")
	}
	return c.Redirect(http.StatusFound, h.frontendRedirect("/"))
}

func (h *AuthHandler) frontendRedirect(path string) string {
	if h.frontendURL == "" {
		return path
	}
	parsed, err := url.Parse(path)
	if err != nil || parsed.IsAbs() || parsed.Host != "" || !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return h.frontendURL + "/"
	}
	return h.frontendURL + path
}

func safeBackPath(path string) string {
	if path == "" || path[0] != '/' || (len(path) > 1 && path[1] == '/') {
		return "/home"
	}
	return path
}

func (h *AuthHandler) CurrentUser(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "ログインが必要です")
	}
	user, err := h.authUsecase.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "ユーザーが見つかりません")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) LinkConnpass(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "認証情報がありません")
	}

	var req struct {
		ConnpassID string `json:"connpass_id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.authUsecase.LinkConnpass(c.Request().Context(), userID, req.ConnpassID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
