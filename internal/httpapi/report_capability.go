package httpapi

import (
	"net/http"

	"github.com/sanskarIN/stockpilot/internal/reporting"
)

func parseReportOffsetParameter(r *http.Request) (int, error) {
	return parseReportOffset(r.URL.Query().Get("offset"))
}

func makeBoundedReportQuery(period reporting.Period, bounds reporting.Bounds) reporting.Query {
	return reporting.Query{
		From:   &period.From,
		To:     &period.To,
		Limit:  bounds.Limit,
		Offset: bounds.Offset,
	}
}
