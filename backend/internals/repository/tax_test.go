package repository

import "testing"

func TestTaxCentsFor(t *testing.T) {
	tests := []struct {
		name    string
		taxable int64
		rateBps int
		want    int64
	}{
		{name: "zero taxable", taxable: 0, rateBps: 1800, want: 0},
		{name: "18 percent", taxable: 10000, rateBps: 1800, want: 1800},
		{name: "rounds half up", taxable: 1999, rateBps: 1800, want: 360},
		{name: "half cent rounds up", taxable: 1, rateBps: 5000, want: 1},
		{name: "below half cent rounds down", taxable: 1, rateBps: 4999, want: 0},
		{name: "zero rate", taxable: 1000, rateBps: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := taxCentsFor(tt.taxable, tt.rateBps); got != tt.want {
				t.Fatalf("taxCentsFor(%d, %d) = %d, want %d", tt.taxable, tt.rateBps, got, tt.want)
			}
		})
	}
}
