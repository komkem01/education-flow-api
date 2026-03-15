package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const maxAuditPayloadLength = 16384

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
	body *bytes.Buffer
}

func (w *auditBodyWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *auditBodyWriter) WriteString(s string) (int, error) {
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func AuditLogMiddleware(auditEnt entitiesinf.AuditLogEntity) gin.HandlerFunc {
	if auditEnt == nil {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		requestBody := captureRequestBody(ctx)

		writer := &auditBodyWriter{ResponseWriter: ctx.Writer, body: bytes.NewBuffer(nil)}
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
			RequestBody:  strPtr(sanitizePayload(requestBody, ctx.ContentType())),
			ResponseBody: strPtr(sanitizePayload(writer.body.String(), ctx.Writer.Header().Get("Content-Type"))),
			ErrorMessage: strPtr(ctx.Errors.String()),
		}

		if fullPath := ctx.FullPath(); fullPath != "" {
			auditLog.RoutePath = &fullPath
		}

		if currentUser, ok := auth.CurrentUserFromGin(ctx); ok && currentUser != nil && currentUser.Member != nil {
			auditLog.ActorID = &currentUser.Member.ID
			role := string(currentUser.Member.Role)
			auditLog.ActorRole = &role
		}

		if spanCtx := trace.SpanContextFromContext(ctx.Request.Context()); spanCtx.IsValid() {
			traceID := spanCtx.TraceID().String()
			auditLog.TraceID = &traceID
		}

		writeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = auditEnt.CreateAuditLog(writeCtx, auditLog)
	}
}

func captureRequestBody(ctx *gin.Context) string {
	if ctx.Request == nil || ctx.Request.Body == nil {
		return ""
	}
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return ""
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

func sanitizePayload(raw string, contentType string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			maskSensitive(v)
			if b, err := json.Marshal(v); err == nil {
				raw = string(b)
			}
		}
	}

	if len(raw) > maxAuditPayloadLength {
		return raw[:maxAuditPayloadLength] + "...(truncated)"
	}
	return raw
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
