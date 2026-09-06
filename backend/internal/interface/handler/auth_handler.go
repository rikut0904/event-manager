package handler

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	pendingMu   sync.Mutex
	pending     map[string]commonid.Pending
	backPaths   map[string]string
	logoutState map[string]time.Time
}

func NewAuthHandler(u usecase.AuthUsecase, commonIDClient *commonid.Client, frontendOrigin string, appSession *session.Manager) *AuthHandler {
	return &AuthHandler{authUsecase: u, commonID: commonIDClient, appSession: appSession, frontendURL: strings.TrimRight(frontendOrigin, "/"), pending: make(map[string]commonid.Pending), backPaths: make(map[string]string), logoutState: make(map[string]time.Time)}
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
	h.pendingMu.Lock()
	h.pending[pending.State] = pending
	h.backPaths[pending.State] = safeBackPath(c.QueryParam("back_path"))
	h.pendingMu.Unlock()
	return c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) Callback(c echo.Context) error {
	if h.commonID == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Common IDが設定されていません")
	}
	state := c.QueryParam("state")
	h.pendingMu.Lock()
	pending, ok := h.pending[state]
	delete(h.pending, state)
	backPath := h.backPaths[state]
	delete(h.backPaths, state)
	h.pendingMu.Unlock()
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "認証状態が見つかりません。もう一度お試しください")
	}
	commonUser, err := h.commonID.Exchange(c.Request().Context(), c.QueryParams(), pending)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if _, err := h.authUsecase.SyncCommonUser(c.Request().Context(), commonUser.CommonUserID, commonUser.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザー情報を保存できませんでした")
	}
	appCookie, err := h.appSession.Issue(commonUser.CommonUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "アプリセッションを発行できませんでした")
	}
	c.SetCookie(appCookie)
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
		return c.Redirect(http.StatusFound, "/")
	}
	logoutURL, state, err := h.commonID.BeginLogout()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ログアウトを開始できません")
	}
	h.pendingMu.Lock()
	h.logoutState[state] = time.Now().Add(10 * time.Minute)
	h.pendingMu.Unlock()
	return c.Redirect(http.StatusFound, logoutURL)
}

func (h *AuthHandler) LogoutCallback(c echo.Context) error {
	c.SetCookie(h.appSession.ClearCookie())
	state := c.QueryParam("state")
	h.pendingMu.Lock()
	expiresAt, ok := h.logoutState[state]
	delete(h.logoutState, state)
	h.pendingMu.Unlock()
	if !ok || time.Now().After(expiresAt) || (c.QueryParam("logout") != "success" && c.QueryParam("result") != "success") {
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
