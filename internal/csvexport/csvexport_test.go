package csvexport

import (
	"bytes"
	"testing"
)

func TestWriterWritesRFC4180CSV(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, Options{})
	if err := w.WriteHeader("sku", "name", "note"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteRow("SKU-1", "Widget, Small", "line one\nline two"); err != nil {
		t.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	const want = "sku,name,note\nSKU-1,\"Widget, Small\",\"line one\nline two\"\n"
	if out.String() != want {
		t.Fatalf("unexpected CSV: %q", out.String())
	}
}

func TestFormulaSafeEscapesSpreadsheetFormulas(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, Options{FormulaSafe: true})
	if err := w.WriteRow("=SUM(A1)", " +1", "-10", "@user", "safe"); err != nil {
		t.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	const want = "'=SUM(A1),\" '+1\",'-10,'@user,safe\n"
	if out.String() != want {
		t.Fatalf("unexpected formula-safe CSV: %q", out.String())
	}
}

func TestWriteHeaderRejectsEmptyHeader(t *testing.T) {
	var out bytes.Buffer
	w := New(&out, Options{})
	if err := w.WriteHeader(); err == nil {
		t.Fatal("expected empty-header error")
	}
}

func TestNilWriterReturnsError(t *testing.T) {
	var w *Writer
	if err := w.WriteRow("x"); err == nil {
		t.Fatal("expected nil writer error")
	}
	if err := w.Flush(); err == nil {
		t.Fatal("expected nil writer flush error")
	}
}
