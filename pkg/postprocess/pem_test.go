package postprocess

import "testing"

func TestIsPem(t *testing.T) {
	tests := []struct {
		name     string
		text     []byte
		isSecret bool
	}{
		{
			"Invalid private pem file",
			[]byte("Random Binary Data"),
			false,
		},
		{
			"Invalid private pem file",
			[]byte("This is not a pem file."),
			false,
		},
		{
			"Invalid private pem file",
			[]byte("Non-pem File"),
			false,
		},
		{
			"Valid private pem file",
			[]byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 69, 78, 67, 82, 89, 80, 84, 69, 68, 32, 80, 82, 73, 86, 65, 84,
				69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 116, 101, 115, 116, 32,
				112, 97, 115, 115, 119, 111, 114, 100, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 69, 78, 67,
				82, 89, 80, 84, 69, 68, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivatePem(tt.text); got != tt.isSecret {
				t.Errorf("TestIsPem() = %v, want %v", got, tt.isSecret)
			}
		})
	}
}
