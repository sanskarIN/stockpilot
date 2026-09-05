package httpapi

import "strings"

// sanitizeCSVCell prefixes spreadsheet formula-leading values so exported text
// is not interpreted as a formula by common spreadsheet applications.
func sanitizeCSVCell(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	switch value[0] {
	case '=', '+', '-', '@':
		return "'" + value
	default:
		return value
	}
}
