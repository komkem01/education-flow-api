package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

func (c *Controller) accessTokenFromRequest(ctx *gin.Context) string {
	if token := extractBearerToken(ctx.GetHeader("Authorization")); token != "" {
		return token
	}

	cookieToken, _ := ctx.Cookie(accessTokenCookieName)
	return strings.TrimSpace(cookieToken)
}

func (c *Controller) refreshTokenFromRequest(ctx *gin.Context, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}

	cookieToken, _ := ctx.Cookie(refreshTokenCookieName)
	return strings.TrimSpace(cookieToken)
}

func (c *Controller) writeAuthCookies(ctx *gin.Context, token string, refreshToken string) {
	ctx.SetSameSite(http.SameSiteLaxMode)

	accessMaxAge := c.svc.conf.AccessTokenTTLSeconds
	refreshMaxAge := c.svc.conf.RefreshTokenTTLSeconds
	if accessMaxAge <= 0 {
		accessMaxAge = 900
	}
	if refreshMaxAge <= 0 {
		refreshMaxAge = 604800
	}

	domain := strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN"))
	secure := ctx.Request.TLS != nil || strings.EqualFold(strings.TrimSpace(os.Getenv("AUTH_COOKIE_SECURE")), "true")

	ctx.SetCookie(accessTokenCookieName, token, accessMaxAge, "/", domain, secure, true)
	ctx.SetCookie(refreshTokenCookieName, refreshToken, refreshMaxAge, "/", domain, secure, true)
}

func (c *Controller) clearAuthCookies(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)

	domain := strings.TrimSpace(os.Getenv("AUTH_COOKIE_DOMAIN"))
	secure := ctx.Request.TLS != nil || strings.EqualFold(strings.TrimSpace(os.Getenv("AUTH_COOKIE_SECURE")), "true")

	ctx.SetCookie(accessTokenCookieName, "", -1, "/", domain, secure, true)
	ctx.SetCookie(refreshTokenCookieName, "", -1, "/", domain, secure, true)
}
