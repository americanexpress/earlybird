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
		{
			name: "Valid private gpg file",
			pgpData: []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 80, 71, 80, 32, 80, 82, 73, 86, 65, 84, 69,
				32, 75, 69, 89, 32, 66, 76, 79, 67, 75, 45, 45, 45, 45, 45, 10, 86, 101, 114, 115, 105, 111, 110,
				58, 32, 66, 67, 80, 71, 32, 118, 49, 46, 54, 48, 10, 10, 115, 111, 109, 101, 100, 117, 109, 121,
				15, 101, 99, 114, 101, 116, 107, 101, 121, 115, 111, 109, 101, 100, 117, 109, 121, 115, 101,
				99, 114, 101, 116, 107, 101, 121, 115, 111, 109, 101, 100, 117, 109, 121, 115, 101, 99, 114,
				101, 116, 107, 101, 121, 115, 111, 109, 101, 100, 117, 109, 121, 115, 101, 99, 114, 101, 116,
				107, 101, 121, 115, 111, 109, 101, 100, 117, 109, 121, 115, 101, 99, 114, 101, 116, 107, 101,
				121, 10, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 80, 71, 80, 32, 80, 82, 73, 86, 65, 84, 69, 32,
				75, 69, 89, 32, 66, 76, 79, 67, 75, 45, 45, 45, 45, 45},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivatePGP(tt.pgpData); got != tt.expected {
				t.Errorf("TestIsPrivatePGP() = %v, want %v", got, tt.expected)
			}
		})
	}
}
