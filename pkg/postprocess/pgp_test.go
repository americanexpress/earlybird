package postprocess

import "testing"

func TestIsPrivatePGP(t *testing.T) {
	tests := []struct {
		name     string
		pgpData  []byte
		expected bool
	}{
		{
			name:     "Invalid GPG file",
			pgpData:  []byte("Random Binary Data"),
			expected: false,
		},
		{
			name:     "Invalid GPG file",
			pgpData:  []byte("This is not a GPG file."),
			expected: false,
		},
		{
			name:     "Invalid GPG file",
			pgpData:  []byte("Non-GPG File"),
			expected: false,
		},
		// {
		// 	name:     "Valid private gpg file",
		// 	pgpData:  []byte{},
		// 	expected: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivatePGP(tt.pgpData); got != tt.expected {
				t.Errorf("TestIsPrivatePGP() = %v, want %v", got, tt.expected)
			}
		})
	}
}
