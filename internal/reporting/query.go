package reporting

import "time"

// Query describes the bounded portion of a report requested by a caller.
// From and To are optional for reports that are not time-based.
type Query struct {
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}
