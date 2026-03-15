package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuditLog struct {
	bun.BaseModel `bun:"table:audit_logs,alias:al"`

	ID           uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	ActorID      *uuid.UUID `bun:"actor_id,type:uuid"`
	ActorRole    *string    `bun:"actor_role,type:varchar(50)"`
	Method       string     `bun:"method,notnull,type:varchar(10)"`
	Path         string     `bun:"path,notnull"`
	RoutePath    *string    `bun:"route_path"`
	StatusCode   int        `bun:"status_code,notnull"`
	LatencyMS    int64      `bun:"latency_ms,notnull"`
	IP           *string    `bun:"ip,type:varchar(64)"`
	UserAgent    *string    `bun:"user_agent"`
	QueryString  *string    `bun:"query_string"`
	RequestBody  *string    `bun:"request_body"`
	ResponseBody *string    `bun:"response_body"`
	ErrorMessage *string    `bun:"error_message"`
	TraceID      *string    `bun:"trace_id,type:varchar(64)"`
	CreatedAt    time.Time  `bun:"created_at,notnull,default:current_timestamp"`
}
