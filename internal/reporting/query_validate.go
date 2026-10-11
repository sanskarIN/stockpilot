package reporting

import "fmt"

// Validate checks repository-facing query bounds without imposing report-specific defaults.
func (q Query) Validate() error {
	if q.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}
	if q.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	if q.From != nil && q.To != nil && q.From.After(*q.To) {
		return fmt.Errorf("from must not be after to")
	}
	return nil
}
