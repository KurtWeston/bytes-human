package converter

import (
	"math"
	"testing"
)

func TestBytesToHuman(t *testing.T) {
	tests := []struct {
		name    string
		bytes   uint64
		opts    Options
		want    string
		wantErr bool
	}{
		{"zero bytes", 0, Options{Standard: Binary, Precision: 1}, "0 B", false},
		{"1024 bytes binary", 1024, Options{Standard: Binary, Precision: 1}, "1.0 KiB", false},
		{"1000 bytes decimal", 1000, Options{Standard: Decimal, Precision: 1}, "1.0 KB", false},
		{"1 MiB", 1048576, Options{Standard: Binary, Precision: 2}, "1.00 MiB", false},
		{"1.5 GB decimal", 1500000000, Options{Standard: Decimal, Precision: 1}, "1.5 GB", false},
		{"force unit KiB", 2048, Options{Standard: Binary, Precision: 0, ForceUnit: "KiB"}, "2 KiB", false},
		{"max uint64", math.MaxUint64, Options{Standard: Binary, Precision: 1}, "16.0 EiB", false},
		{"invalid precision", 1024, Options{Standard: Binary, Precision: 7}, "", true},
		{"invalid unit", 1024, Options{Standard: Binary, Precision: 1, ForceUnit: "XB"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BytesToHuman(tt.bytes, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesToHuman() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BytesToHuman() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanToBytes(t *testing.T) {
	tests := []struct {
		name    string
		human   string
		opts    Options
		want    uint64
		wantErr bool
	}{
		{"1 KiB", "1 KiB", Options{Standard: Binary}, 1024, false},
		{"1.5 GB", "1.5 GB", Options{Standard: Decimal}, 1500000000, false},
		{"1024 B", "1024 B", Options{Standard: Binary}, 1024, false},
		{"2.5 MiB", "2.5 MiB", Options{Standard: Binary}, 2621440, false},
		{"no space", "1KB", Options{Standard: Decimal}, 1000, false},
		{"lowercase", "1kb", Options{Standard: Decimal}, 1000, false},
		{"invalid format", "abc", Options{Standard: Binary}, 0, true},
		{"negative number", "-1 KB", Options{Standard: Decimal}, 0, true},
		{"unknown unit", "1 XB", Options{Standard: Binary}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HumanToBytes(tt.human, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("HumanToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HumanToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundingModes(t *testing.T) {
	tests := []struct {
		name  string
		bytes uint64
		mode  RoundMode
		want  string
	}{
		{"round nearest", 1536, RoundNearest, "1.5 KiB"},
		{"round up", 1536, RoundUp, "1.5 KiB"},
		{"round down", 1536, RoundDown, "1.5 KiB"},
		{"round up fraction", 1540, RoundUp, "1.6 KiB"},
		{"round down fraction", 1540, RoundDown, "1.5 KiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{Standard: Binary, Precision: 1, RoundMode: tt.mode}
			got, err := BytesToHuman(tt.bytes, opts)
			if err != nil {
				t.Errorf("BytesToHuman() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("BytesToHuman() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyRounding(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		precision int
		mode      RoundMode
		want      float64
	}{
		{"nearest", 1.555, 2, RoundNearest, 1.56},
		{"up", 1.551, 2, RoundUp, 1.56},
		{"down", 1.559, 2, RoundDown, 1.55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyRounding(tt.value, tt.precision, tt.mode)
			if got != tt.want {
				t.Errorf("applyRounding() = %v, want %v", got, tt.want)
			}
		})
	}
}
