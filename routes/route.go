package routes

import (
	"net/http"
	"os"
	"strings"

	"eduflow/app/modules"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func Router(app *gin.Engine, mod *modules.Modules) {
	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	app.Use(otelgin.Middleware(mod.Conf.Svc.Config().AppName),
		// Middleware add trace id to response header
		func(ctx *gin.Context) {
			spanCtx := trace.SpanContextFromContext(ctx.Request.Context())
			if spanCtx.IsValid() {
				ctx.Header("X-Trace-ID", spanCtx.TraceID().String())
			}
			ctx.Next()
		},
	)

	origins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	if raw := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS")); raw != "" {
		origins = make([]string, 0, len(strings.Split(raw, ",")))
		for _, item := range strings.Split(raw, ",") {
			origin := strings.TrimSpace(item)
			if origin != "" {
				origins = append(origins, origin)
			}
		}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:           origins,
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:          []string{"X-Trace-ID"},
		AllowCredentials:       true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}))

	v1 := app.Group("/api/v1")
	v1.Use(AuditLogMiddleware(mod.ENT.Svc))

	api(v1, mod)
	apiSystem(v1, mod)
	apiPublic(v1, mod)
	apiMember(v1, mod)
	apiManagement(v1, mod)
	apiEnrollments(v1, mod)
	apiAttendance(v1, mod)
	apiApproval(v1, mod)
}
