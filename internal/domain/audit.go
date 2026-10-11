package domain

import (
	"fmt"
	"strings"
	"time"
)

func (e AuditEvent) Validate() error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now().UTC()
	}
	if strings.TrimSpace(e.Action) == "" || len(strings.TrimSpace(e.Action)) > 80 {
		return fmt.Errorf("%w: audit action must be 1-80 characters", ErrInvalid)
	}
	if strings.TrimSpace(e.EntityType) == "" || len(strings.TrimSpace(e.EntityType)) > 80 {
		return fmt.Errorf("%w: audit entity type must be 1-80 characters", ErrInvalid)
	}
	if strings.TrimSpace(e.EntityID) == "" || len(strings.TrimSpace(e.EntityID)) > 256 {
		return fmt.Errorf("%w: audit entity id must be 1-256 characters", ErrInvalid)
	}
	if e.Metadata == nil {
		e.Metadata = map[string]any{}
	}
	return nil
}
