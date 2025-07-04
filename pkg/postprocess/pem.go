package postprocess

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strings"
)

func IsPrivatePem(filePath string) bool {
	// Read the PEM file
	pemData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read PEM file %v: %v", filePath, err)
		return true
	}

	// Loop through all PEM blocks
	for {
		block, rest := pem.Decode(pemData)
		if block == nil {
			return false // No more blocks to decode
		}

		if checkPrivate(block.Type) {
			return true
		}
		if checkPublic(block.Type) {
			pemData = rest
			continue
		}

		// Attempt to parse as a private key
		_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			return true
		}

		// Attempt to parse as a public key
		_, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err == nil {
			pemData = rest
			continue
		}

		// Update pemData to process the remaining blocks
		pemData = rest
	}
}

func checkPrivate(blockType string) bool {
	list := []string{"private", "PKCS8"}
	joinedList := strings.Join(list, " ")
	blockTypes := strings.Split(strings.ToLower(blockType), " ")

	for _, bt := range blockTypes {
		if strings.Contains(joinedList, bt) {
			return true
		}
	}
	return false
}

func checkPublic(blockType string) bool {
	list := []string{"public", "certificate", "x509"}
	joinedList := strings.Join(list, " ")
	blockTypes := strings.Split(strings.ToLower(blockType), " ")

	for _, bt := range blockTypes {
		if strings.Contains(joinedList, bt) {
			return true
		}
	}
	return false
}
