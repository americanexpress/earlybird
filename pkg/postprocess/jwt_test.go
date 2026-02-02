package postprocess

import "testing"

func TestProcessJWT(t *testing.T) {
	// Your test code here
	tests := []struct {
		name     string
		text     string
		isSecret bool
	}{
		{
			"Valid JWT token, as it does not contain expiration claim",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQW1lcmljYW4gRXhwcmVzcyIsImFkbWluIjp0cnVlLCJpYXQiOjE1MTYyMzkwMjJ9.JieTvCL1z9dCSAoFAtYWkHEKswbfEQGXCVcej2g7n80",
			true,
		},
		{
			"Invalid JWT token, as it contains expired expiration claim",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjIsIm5hbWUiOiJBbWVyaWNhbiBFeHByZXNzIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.nTCTRxHQRCbmIplyVwZzPw8_Ohv-3wJpskH_MphhUFE",
			false,
		},
		{
			"Invalid JWT token, malformed token which are failed to parse are considered valid jwt findings",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVJ9.eyJhbGciOiJIUzINiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsInR5cCI6kpXVCJ9",
			true,
		},
		{
			"Ignore finding for Encrypted JWE token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			false,
		},
		{
			"Valid JWT token, as it does contain expiration claim, and the expiration time is in the future",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMxNTU4NjkxNzkzMDAwLCJuYW1lIjoiQW1lcmljYW4gRXhwcmVzcyIsImFkbWluIjp0cnVlLCJpYXQiOjE1MTYyMzkwMjJ9.wOSLG3VTKg7A2KoKu9WsuIiKIpl__PDWStw1CkmeOPs",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidJWT(tt.text); got != tt.isSecret {
				t.Errorf("IsValidJWT() = %v, want %v", got, tt.isSecret)
			}
		})
	}
}
