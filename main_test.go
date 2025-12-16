package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/user/bytes-human/converter"
)

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		cfg     config
		want    string
		wantErr bool
	}{
		{
			name:  "bytes to human",
			input: "1024",
			cfg:   config{toHuman: true, precision: 1},
			want:  "1.0 KiB",
		},
		{
			name:  "human to bytes",
			input: "1.5 GB",
			cfg:   config{toBytes: true, decimal: true},
			want:  "1500000000",
		},
		{
			name:  "decimal units",
			input: "1000000",
			cfg:   config{toHuman: true, decimal: true, precision: 1},
			want:  "1.0 MB",
		},
		{
			name:    "invalid bytes",
			input:   "abc",
			cfg:     config{toHuman: true},
			wantErr: true,
		},
		{
			name:    "invalid human format",
			input:   "invalid",
			cfg:     config{toBytes: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processInput(tt.input, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("processInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("processInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessBatch(t *testing.T) {
	input := "1024\n2048\n4096\n"
	old := os.Stdin
	defer func() { os.Stdin = old }()

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	var buf bytes.Buffer
	old2 := os.Stdout
	os.Stdout = nil
	defer func() { os.Stdout = old2 }()

	cfg := config{toHuman: true, precision: 1, batch: true}
	processBatch(cfg)
}

func TestParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-h", "-p", "2"}
	cfg := parseFlags()

	if !cfg.toHuman {
		t.Error("Expected toHuman to be true")
	}
	if cfg.precision != 2 {
		t.Errorf("Expected precision 2, got %d", cfg.precision)
	}
}

func TestResultJSON(t *testing.T) {
	r := result{
		Input:  "1024",
		Output: "1.0 KiB",
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	if !strings.Contains(string(data), "1024") {
		t.Error("JSON output missing input")
	}
	if !strings.Contains(string(data), "1.0 KiB") {
		t.Error("JSON output missing output")
	}
}
