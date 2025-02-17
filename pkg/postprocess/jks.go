package postprocess

import (
	"github.com/americanexpress/earlybird/v4/pkg/jks"
)

// JKS checks if a file meets standard
func JKS(fileByte []byte) bool {
	ks, err := jks.Parse(fileByte)
	if err != nil {
		println("JKS Parse Error", err)
	}
	if ks != nil && len(ks.Keypairs) > 0 {
		return true
	}
	return false
}
