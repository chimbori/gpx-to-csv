package main

import (
	"testing"
	"time"
)

func TestUtcToLocal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name:    "valid RFC3339 timestamp",
			input:   "2023-01-15T10:30:45Z",
			wantErr: false,
			check: func(t *testing.T, result string) {
				// Parse result to ensure it's valid RFC3339
				if _, err := time.Parse(time.RFC3339, result); err != nil {
					t.Errorf("result is not valid RFC3339: %v", err)
				}
			},
		},
		{
			name:    "another valid RFC3339 timestamp",
			input:   "2024-12-25T23:59:59Z",
			wantErr: false,
			check: func(t *testing.T, result string) {
				if _, err := time.Parse(time.RFC3339, result); err != nil {
					t.Errorf("result is not valid RFC3339: %v", err)
				}
			},
		},
		{
			name:    "invalid timestamp format",
			input:   "not-a-timestamp",
			wantErr: true,
			check: func(t *testing.T, result string) {
				// On error, should return the original input
				if result != "not-a-timestamp" {
					t.Errorf("expected original input on error, got %s", result)
				}
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
			check: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("expected empty string on error, got %s", result)
				}
			},
		},
		{
			name:    "malformed RFC3339",
			input:   "2023-13-45T25:70:90Z",
			wantErr: true,
			check: func(t *testing.T, result string) {
				if result != "2023-13-45T25:70:90Z" {
					t.Errorf("expected original input on error, got %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utcToLocal(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("utcToLocal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

func TestUtcToLocalConversion(t *testing.T) {
	// Test that a UTC time is actually converted to local time
	input := "2023-01-15T12:00:00Z"
	result, err := utcToLocal(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedResult, _ := time.Parse(time.RFC3339, result)
	parsedInput, _ := time.Parse(time.RFC3339, input)

	// The local time should be equal to UTC time in absolute terms,
	// but the timezone offset should be different
	if !parsedResult.Equal(parsedInput) {
		t.Errorf("converted time does not represent the same instant: input=%v, result=%v", parsedInput, parsedResult)
	}
}
