package postprocess

import "testing"

func TestIsJKS(t *testing.T) {
	type args struct {
		cc string
	}
	tests := []struct {
		name     string
		text     []byte
		isSecret bool
	}{
		{
			"Invalid JKS file",
			[]byte("Random Binary Data"),
			false,
		},
		{
			"Valid Basic Auth Header",
			[]byte("This is not a JKS file."),
			false,
		},
		{
			"Empty Basic Auth Header",
			[]byte("Non-JKS File"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JKS(tt.text); got != tt.isSecret {
				t.Errorf("TestIsJKS() = %v, want %v", got, tt.isSecret)
			}
		})
	}
}
