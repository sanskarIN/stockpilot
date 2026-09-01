package httpapi

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (a *API) recordAudit(ctx context.Context, actorID, action, entityType, entityID string, metadata map[string]any) {
	if a.audit == nil {
		return
	}
	event := domain.AuditEvent{
		OccurredAt: time.Now().UTC(),
		ActorID: actorID,
		Action: action,
		EntityType: entityType,
		EntityID: entityID,
		RequestID: requestIDFromContext(ctx),
		Metadata: metadata,
	}
	if err := a.audit.AppendAuditEvent(ctx, event); err != nil && a.logger != nil {
		a.logger.ErrorContext(ctx, "audit append failed", slog.Any("error", err), "action", action, "entity_type", entityType, "entity_id", entityID)
	}
}

func requestIDFromContext(ctx context.Context) string {
	if value, ok := ctx.Value(requestIDContextKey{}).(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

type requestIDContextKey struct{}
