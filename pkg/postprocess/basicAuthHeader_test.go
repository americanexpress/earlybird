package postprocess

import "testing"

func TestIsBasicAuthHeader(t *testing.T) {
	type args struct {
		cc string
	}
	tests := []struct {
		name     string
		header   string
		isSecret bool
	}{
		{
			"Valid Basic Auth Header",
			"Authorization: Basic dXNlbmFtZTpwYXNzd29yZA==",
			true,
		},
		{
			"Invalid Basic Auth Header",
			"Authorization: Basic acdefjhikl",
			false,
		},
		{
			"garbage in middle - this is handled by the pattern matching",
			"370000ajlsdklasdj000000002",
			true,
		},
		{
			"Test AMEX example credit card",
			"370000000000002",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBasicAuthHeader(tt.header); got != tt.isSecret {
				t.Errorf("IsBasicAuthHeader() = %v, want %v", got, tt.isSecret)
			}
		})
	}
}
