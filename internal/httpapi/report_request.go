package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/reporting"
)

const (
	defaultReportLimit = 100
	maxReportLimit     = 5000
)

// reportRequest captures the shared, validated query contract for reporting endpoints.
type reportRequest struct {
	Period reporting.Period
	Bounds reporting.Bounds
}

func parseReportRequest(r *http.Request, now time.Time) (reportRequest, error) {
	q := r.URL.Query()
	fromText, toText := q.Get("from"), q.Get("to")

	var from, to time.Time
	var err error
	if fromText == "" && toText == "" {
		to = now.UTC()
		from = to.AddDate(0, 0, -(reporting.DefaultPeriodDays - 1))
	} else if fromText == "" || toText == "" {
		return reportRequest{}, fmt.Errorf("reporting period requires both from and to")
	} else {
		from, err = time.Parse(time.RFC3339, fromText)
		if err != nil {
			return reportRequest{}, fmt.Errorf("invalid from timestamp")
		}
		to, err = time.Parse(time.RFC3339, toText)
		if err != nil {
			return reportRequest{}, fmt.Errorf("invalid to timestamp")
		}
	}

	period, err := reporting.NewPeriod(from, to)
	if err != nil {
		return reportRequest{}, err
	}

	limit, err := parseReportLimit(q.Get("limit"))
	if err != nil {
		return reportRequest{}, err
	}
	offset, err := parseReportOffset(q.Get("offset"))
	if err != nil {
		return reportRequest{}, err
	}

	return reportRequest{Period: period, Bounds: reporting.NewBounds(limit, offset)}, nil
}

func parseReportLimit(value string) (int, error) {
	if value == "" {
		return defaultReportLimit, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("limit must be a positive integer")
	}
	if parsed > maxReportLimit {
		return 0, fmt.Errorf("limit must not exceed %d", maxReportLimit)
	}
	return parsed, nil
}

func parseReportOffset(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("offset must be a non-negative integer")
	}
	return parsed, nil
}

func writeReportMetadata(w http.ResponseWriter, request reportRequest, generatedAt time.Time, complete bool) {
	metadata := reporting.NewMetadata(request.Period, request.Bounds, complete, generatedAt)
	w.Header().Set("X-Report-Generated-At", metadata.GeneratedAt.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-From", metadata.From.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-To", metadata.To.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-Limit", strconv.Itoa(metadata.Limit))
	w.Header().Set("X-Report-Offset", strconv.Itoa(metadata.Offset))
	w.Header().Set("X-Report-Complete", strconv.FormatBool(metadata.Complete))
}
