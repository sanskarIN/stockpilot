package csvexport

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Options controls CSV serialization behavior.
type Options struct {
	// FormulaSafe prefixes a single quote to cells that spreadsheet programs may
	// interpret as formulas. It is disabled by default so generic CSV exports
	// preserve their source values.
	FormulaSafe bool
}

// Writer serializes tabular data as RFC 4180-compatible CSV.
type Writer struct {
	w      *csv.Writer
	option Options
}

// New returns a CSV writer that writes to dst.
func New(dst io.Writer, option Options) *Writer {
	return &Writer{w: csv.NewWriter(dst), option: option}
}

// WriteHeader writes the column header row. An empty header is rejected.
func (w *Writer) WriteHeader(columns ...string) error {
	if len(columns) == 0 {
		return fmt.Errorf("csvexport: at least one column is required")
	}
	return w.WriteRow(columns...)
}

// WriteRow writes one row. The row may contain an arbitrary number of fields.
func (w *Writer) WriteRow(fields ...string) error {
	if w == nil || w.w == nil {
		return fmt.Errorf("csvexport: writer is nil")
	}
	row := make([]string, len(fields))
	for i, field := range fields {
		row[i] = sanitize(field, w.option.FormulaSafe)
	}
	if err := w.w.Write(row); err != nil {
		return fmt.Errorf("csvexport: write row: %w", err)
	}
	return nil
}

// Flush completes buffered output and reports any underlying write error.
func (w *Writer) Flush() error {
	if w == nil || w.w == nil {
		return fmt.Errorf("csvexport: writer is nil")
	}
	w.w.Flush()
	if err := w.w.Error(); err != nil {
		return fmt.Errorf("csvexport: flush: %w", err)
	}
	return nil
}

func sanitize(value string, formulaSafe bool) string {
	if !formulaSafe || value == "" {
		return value
	}
	trimmed := strings.TrimLeft(value, " \t")
	if trimmed == "" {
		return value
	}
	switch trimmed[0] {
	case '=', '+', '-', '@':
		prefix := value[:len(value)-len(trimmed)]
		return prefix + "'" + trimmed
	default:
		return value
	}
}
