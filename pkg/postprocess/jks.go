package postprocess

import (
	"fmt"
	"github.com/americanexpress/earlybird/v4/pkg/jks"
	"os"
)

// JKS checks if a file meets standard
func JKS(file string) bool {

	raw, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("JKS file read Error")
		return false
	}

	ks, err := jks.Parse(raw)
	if err != nil {
		println("JKS Parse Error", err)
	}
	if ks != nil && len(ks.Keypairs) > 0 {
		return true
	}
	return false
}
