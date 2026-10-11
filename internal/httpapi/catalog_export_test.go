package httpapi

import "testing"

func TestNormalizeCatalogExportBounds(t *testing.T) {
	tests := []struct {
		name          string
		limit, offset int
		wantLimit     int
		wantOffset    int
	}{
		{name: "defaults invalid limit and offset", limit: 0, offset: -4, wantLimit: defaultCatalogExportRows, wantOffset: 0},
		{name: "clamps maximum", limit: maxCatalogExportRows + 500, offset: 8, wantLimit: maxCatalogExportRows, wantOffset: 8},
		{name: "preserves valid values", limit: 250, offset: 12, wantLimit: 250, wantOffset: 12},
		{name: "negative limit uses default", limit: -10, offset: 0, wantLimit: defaultCatalogExportRows, wantOffset: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := normalizeCatalogExportBounds(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
				t.Fatalf("normalizeCatalogExportBounds(%d, %d) = (%d, %d), want (%d, %d)", tt.limit, tt.offset, gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}
