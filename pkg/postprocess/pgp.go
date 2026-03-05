package postprocess

import (
	"bytes"

	"github.com/ProtonMail/go-crypto/openpgp"
)

// IsPrivatePGP checks if the provided data contains a PGP private key.
func IsPrivatePGP(pgpData []byte) bool {
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(pgpData))
	if err != nil {
		return false
	}
	for _, entity := range entityList {
		if entity.PrivateKey != nil {
			return true
		}
		// Also check for subkeys
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil {
				return true
			}
		}
	}
	return false
}
