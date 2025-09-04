package postprocess

import "testing"

func TestIsPem(t *testing.T) {
	type args struct {
		cc string
	}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivatePem(tt.text); got != tt.isSecret {
				t.Errorf("TestIsPem() = %v, want %v", got, tt.isSecret)
			}
		})
	}
}
