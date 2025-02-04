package jks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"testing"
)

// TestOIDFromNamedCurve ensures that we return the correct OID identifying the
// curve for ECDSA private keys.
func TestOIDFromNamedCurve(t *testing.T) {
	t.Run("P-224", testOIDFromNamedCurve(oidNamedCurveP224, elliptic.P224()))
	t.Run("P-256", testOIDFromNamedCurve(oidNamedCurveP256, elliptic.P256()))
	t.Run("P-384", testOIDFromNamedCurve(oidNamedCurveP384, elliptic.P384()))
	t.Run("P-521", testOIDFromNamedCurve(oidNamedCurveP521, elliptic.P521()))
}

func testOIDFromNamedCurve(exp asn1.ObjectIdentifier, curve elliptic.Curve,
) func(*testing.T) {
	return func(t *testing.T) {
		k, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			t.Fatalf("failed to generate key: %v", err)
		}

		oid, err := oidFromNamedCurve(k)
		switch {
		case err != nil:
			t.Errorf("could not find OID: %v", err)
		case !oid.Equal(exp):
			t.Errorf("OID %v ≠ expected %v", oid, exp)
		}
	}
}
