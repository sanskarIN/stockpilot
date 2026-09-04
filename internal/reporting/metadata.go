package reporting

import "time"

// Metadata describes the execution context returned alongside a bounded report.
type Metadata struct {
	GeneratedAt time.Time `json:"generatedAt"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	Limit       int       `json:"limit"`
	Offset      int       `json:"offset"`
	Complete    bool      `json:"complete"`
}

func NewMetadata(period Period, bounds Bounds, complete bool, generatedAt time.Time) Metadata {
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	return Metadata{
		GeneratedAt: generatedAt,
		From:        period.From,
		To:          period.To,
		Limit:       bounds.Limit,
		Offset:      bounds.Offset,
		Complete:    complete,
	}
}
