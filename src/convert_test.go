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

func TestPrecision7digit(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "positive latitude",
			input:    40.7128,
			expected: "40.7128000",
		},
		{
			name:     "negative latitude",
			input:    -74.0060,
			expected: "-74.0060000",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0.0000000",
		},
		{
			name:     "small positive number",
			input:    0.0000001,
			expected: "0.0000001",
		},
		{
			name:     "small negative number",
			input:    -0.0000001,
			expected: "-0.0000001",
		},
		{
			name:     "number with more than 7 decimal places",
			input:    3.141592653589793,
			expected: "3.1415927",
		},
		{
			name:     "large positive number",
			input:    123456.789,
			expected: "123456.7890000",
		},
		{
			name:     "large negative number",
			input:    -123456.789,
			expected: "-123456.7890000",
		},
		{
			name:     "single decimal place",
			input:    45.5,
			expected: "45.5000000",
		},
		{
			name:     "integer-like float",
			input:    90.0,
			expected: "90.0000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := precision7digit(tt.input)
			if result != tt.expected {
				t.Errorf("precision7digit(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPrecision7digitDecimalCount(t *testing.T) {
	// Verify that the result always has exactly 7 decimal places
	testValues := []float64{1.0, -1.0, 0.0, 123.456, -789.012, 3.141592653589793}

	for _, val := range testValues {
		result := precision7digit(val)

		// Find the decimal point
		decimalIndex := -1
		for i, ch := range result {
			if ch == '.' {
				decimalIndex = i
				break
			}
		}

		if decimalIndex == -1 {
			t.Errorf("precision7digit(%v) = %s, expected to contain a decimal point", val, result)
			continue
		}

		// Count decimal places
		decimalPlaces := len(result) - decimalIndex - 1
		if decimalPlaces != 7 {
			t.Errorf("precision7digit(%v) = %s has %d decimal places, expected 7", val, result, decimalPlaces)
		}
	}
}
