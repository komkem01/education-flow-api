package routes

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const maxAuditPayloadLength = 16384
const maxAuditReadBytes = 16384
const auditRetentionDays = 90
const auditPurgeInterval = 6 * time.Hour

var auditPurgeMu sync.Mutex
var auditLastPurgeAt time.Time

var auditSensitiveKeys = map[string]struct{}{
	"password":      {},
	"token":         {},
	"refresh_token": {},
	"access_token":  {},
	"authorization": {},
	"secret":        {},
	"jwt":           {},
}

type auditBodyWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	max       int
	truncated bool
}

func (w *auditBodyWriter) Write(b []byte) (int, error) {
	if w.body != nil && w.max > 0 {
		remain := w.max - w.body.Len()
		if remain > 0 {
			if len(b) > remain {
				w.body.Write(b[:remain])
				w.truncated = true
			} else {
				w.body.Write(b)
			}
		} else {
			w.truncated = true
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *auditBodyWriter) WriteString(s string) (int, error) {
	if w.body != nil && w.max > 0 {
		remain := w.max - w.body.Len()
		if remain > 0 {
			if len(s) > remain {
				w.body.WriteString(s[:remain])
				w.truncated = true
			} else {
				w.body.WriteString(s)
			}
		} else {
			w.truncated = true
		}
	}
	return w.ResponseWriter.WriteString(s)
}

func (w *auditBodyWriter) BodyString() string {
	if w.body == nil {
		return ""
	}
	raw := w.body.String()
	if w.truncated {
		return raw + "...(truncated)"
	}
	return raw
}

func AuditLogMiddleware(auditEnt entitiesinf.AuditLogEntity) gin.HandlerFunc {
	if auditEnt == nil {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		requestBody := captureRequestBody(ctx)

		writer := &auditBodyWriter{ResponseWriter: ctx.Writer, body: bytes.NewBuffer(nil), max: maxAuditReadBytes}
		ctx.Writer = writer
		ctx.Next()

		auditLog := &ent.AuditLog{
			Method:       ctx.Request.Method,
			Path:         ctx.Request.URL.Path,
			StatusCode:   ctx.Writer.Status(),
			LatencyMS:    time.Since(start).Milliseconds(),
			IP:           strPtr(ctx.ClientIP()),
			UserAgent:    strPtr(ctx.Request.UserAgent()),
			QueryString:  strPtr(ctx.Request.URL.RawQuery),
			RequestBody:  strPtr(sanitizePayload(requestBody, ctx.ContentType(), ctx.Request.URL.Path)),
			ResponseBody: strPtr(sanitizePayload(writer.BodyString(), ctx.Writer.Header().Get("Content-Type"), ctx.Request.URL.Path)),
			ErrorMessage: strPtr(ctx.Errors.String()),
		}

		if fullPath := ctx.FullPath(); fullPath != "" {
			auditLog.RoutePath = &fullPath
		}

		if currentUser, ok := auth.CurrentUserFromGin(ctx); ok && currentUser != nil && currentUser.Member != nil {
			auditLog.ActorID = &currentUser.Member.ID
			role := string(currentUser.Member.Role)
			auditLog.ActorRole = &role
		} else if actorID, actorRole := parseActorFromBearer(ctx.GetHeader("Authorization")); actorID != nil {
			auditLog.ActorID = actorID
			auditLog.ActorRole = actorRole
		}

		if spanCtx := trace.SpanContextFromContext(ctx.Request.Context()); spanCtx.IsValid() {
			traceID := spanCtx.TraceID().String()
			auditLog.TraceID = &traceID
		}

		writeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = auditEnt.CreateAuditLog(writeCtx, auditLog)

		go maybePurgeAuditLogs(auditEnt)
	}
}

func captureRequestBody(ctx *gin.Context) string {
	if ctx.Request == nil || ctx.Request.Body == nil {
		return ""
	}
	if !shouldCaptureRequestBody(ctx.Request.Method, ctx.ContentType(), ctx.Request.ContentLength) {
		return ""
	}
	bodyBytes, err := io.ReadAll(io.LimitReader(ctx.Request.Body, maxAuditReadBytes+1))
	if err != nil {
		return ""
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if len(bodyBytes) > maxAuditReadBytes {
		return string(bodyBytes[:maxAuditReadBytes]) + "...(truncated)"
	}
	return string(bodyBytes)
}

func sanitizePayload(raw string, contentType string, path string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	lowerContentType := strings.ToLower(contentType)
	if shouldRedactPath(path) {
		return "[REDACTED]"
	}

	if strings.Contains(lowerContentType, "application/json") {
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			maskSensitive(v)
			if b, err := json.Marshal(v); err == nil {
				raw = string(b)
			}
		}
	}
	if strings.Contains(lowerContentType, "application/x-www-form-urlencoded") {
		if vals, err := url.ParseQuery(raw); err == nil {
			for k := range vals {
				if _, ok := auditSensitiveKeys[strings.ToLower(k)]; ok {
					vals.Set(k, "***")
				}
			}
			raw = vals.Encode()
		}
	}
	if strings.Contains(lowerContentType, "multipart/form-data") {
		return "[MULTIPART_OMITTED]"
	}

	if len(raw) > maxAuditPayloadLength {
		return raw[:maxAuditPayloadLength] + "...(truncated)"
	}
	return raw
}

func shouldCaptureRequestBody(method string, contentType string, contentLength int64) bool {
	switch strings.ToUpper(method) {
	case "GET", "HEAD", "OPTIONS":
		return false
	}
	lowerContentType := strings.ToLower(contentType)
	if strings.Contains(lowerContentType, "multipart/form-data") {
		return false
	}
	if contentLength < 0 || contentLength > maxAuditReadBytes {
		return false
	}
	return true
}

func shouldRedactPath(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "/auth/login") || strings.Contains(p, "/auth/refresh")
}

func maybePurgeAuditLogs(auditEnt entitiesinf.AuditLogEntity) {
	now := time.Now()
	auditPurgeMu.Lock()
	if !auditLastPurgeAt.IsZero() && now.Sub(auditLastPurgeAt) < auditPurgeInterval {
		auditPurgeMu.Unlock()
		return
	}
	auditLastPurgeAt = now
	auditPurgeMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = auditEnt.PurgeAuditLogsBefore(ctx, now.AddDate(0, 0, -auditRetentionDays))
}

func maskSensitive(v any) {
	switch val := v.(type) {
	case map[string]any:
		for key, child := range val {
			if _, ok := auditSensitiveKeys[strings.ToLower(key)]; ok {
				val[key] = "***"
				continue
			}
			maskSensitive(child)
		}
	case []any:
		for _, child := range val {
			maskSensitive(child)
		}
	}
}

func strPtr(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}

type auditAccessClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
}

func parseActorFromBearer(bearer string) (*uuid.UUID, *string) {
	token := extractBearerToken(bearer)
	if token == "" {
		return nil, nil
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, nil
	}

	unsigned := parts[0] + "." + parts[1]
	if !isValidTokenSignature(unsigned, parts[2], secret) {
		return nil, nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil
	}

	claims := new(auditAccessClaims)
	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, nil
	}

	parsedID, err := uuid.Parse(strings.TrimSpace(claims.Sub))
	if err != nil {
		return nil, nil
	}

	var role *string
	if strings.TrimSpace(claims.Role) != "" {
		r := claims.Role
		role = &r
	}

	return &parsedID, role
}

func extractBearerToken(v string) string {
	if v == "" {
		return ""
	}
	parts := strings.SplitN(strings.TrimSpace(v), " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func isValidTokenSignature(unsigned string, sig string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsigned))
	expected := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(sig), []byte(expected))
}
