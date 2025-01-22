/*
Package jks provides routines for manipulating Java Keystore files.
*/
package jks

import (
	"bytes"
	"crypto/hmac"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

// Parse a JKS file. If desired, opts may be specified to provide more control
// over the parsing. If nil, then we will use an empty password when attempting
// to decrypt keys and will not attempt to verify the digest stored in the file.
//
// Errors encountered when parsing a certificate, or decrypting or parsing a
// private key, are stored within the returned Keystore structure. These do not
// lead to the parse failing and will not be returned as an error by the top
// level function. Unrecoverable errors (i.e. malformed file) will result in the
// Parse function returning an error. If digest verification is requested and
// the password or the digest is incorrect, an error will also be returned. If
// any useful data has been extracted it will be returned as a partial Keystore.
func Parse(raw []byte, opts *Options) (*Keystore, error) {
	if opts == nil {
		opts = &defaultOptions
	}

	buf := bytes.NewReader(raw)
	ks := new(Keystore)

	// read file header
	magic, _, err := readUint32(buf, "magic header")
	if err != nil {
		println("err in magic", err)
		return nil, err
	}
	if magic != MagicNumber {
		return nil, fmt.Errorf("invalid magic; expected 0x%08X "+
			"but got 0x%08X", MagicNumber, magic)
	}

	version, _, err := readUint32(buf, "file version")
	if err != nil {
		println("err in version", err)
		return nil, err
	}
	if version != 2 {
		return nil, fmt.Errorf("found version %d file, but expected "+
			"version 2", version)
	}

	numEnts, _, err := readUint32(buf, "number of entries")
	if err != nil {
		println("err in entries", err)

		return nil, err
	}

	// read each entry in turn
	for n := uint32(0); n < numEnts; n++ {
		etype, pos, err := readUint32(buf, "entry type")
		if err != nil {
			println("err in entry type", err)
			return ks, err
		}
		switch etype {
		case 1:
			// it's a private key + cert chain
			kp, err := readKeypair(buf, opts)
			if err != nil {
				return ks, err
			}
			ks.Keypairs = append(ks.Keypairs, kp)

		case 2:
			// it's a certificate
			cert, err := readCert(buf)
			if err != nil {
				return ks, err
			}
			ks.Certs = append(ks.Certs, cert)

		default:
			return nil, fmt.Errorf("unrecognised entry type %d "+
				"at file position %d", etype, pos)
		}
	}

	switch {
	// there should be exactly 20 bytes left
	case buf.Len() != 20:
		return ks, errors.New("malformed digest at end of file")

	case opts.SkipVerifyDigest:
		return ks, nil

	default:
		digest := ComputeDigest(raw[:len(raw)-20], opts.Password)
		if !hmac.Equal(digest, raw[len(raw)-20:]) {
			return ks, errors.New("digest mismatch")
		}
		return ks, nil
	}
}

func readUint32(buf *bytes.Reader, desc string,
) (value uint32, offset int64, err error) {
	offset, _ = buf.Seek(0, io.SeekCurrent)
	if buf.Len() < 4 {
		return 0, offset, fmt.Errorf("unexpected EOF at position %d "+
			"while reading %s", offset, desc)
	}

	var raw [4]byte
	_, _ = buf.Read(raw[:])
	return binary.BigEndian.Uint32(raw[:]), offset, nil
}

func readUint64(buf *bytes.Reader, desc string,
) (value uint64, offset int64, err error) {
	offset, _ = buf.Seek(0, io.SeekCurrent)
	if buf.Len() < 8 {
		return 0, offset, fmt.Errorf("unexpected EOF at position %d "+
			"while reading %s", offset, desc)
	}

	var raw [8]byte
	_, _ = buf.Read(raw[:])
	return binary.BigEndian.Uint64(raw[:]), offset, nil
}

func readTimestamp(buf *bytes.Reader) (ts time.Time, offset int64, err error) {
	ums, offset, err := readUint64(buf, "timestamp")
	if err != nil {
		return time.Time{}, offset, err
	}
	ms := int64(ums)
	return time.Unix(ms/1000, (ms%1000)*1e6), offset, nil
}

func readStr(buf *bytes.Reader, desc string,
) (value string, offset int64, err error) {
	offset, _ = buf.Seek(0, io.SeekCurrent)
	if buf.Len() < 2 {
		return "", offset, fmt.Errorf("unexpected EOF at position %d "+
			"while reading %s", offset, desc)
	}

	var raw [2]byte
	_, _ = buf.Read(raw[:])
	strlen := binary.BigEndian.Uint16(raw[:])
	if buf.Len() < 2 {
		return "", offset, fmt.Errorf("unexpected EOF at position %d "+
			"while reading %s (stored length %d)",
			offset, desc, strlen)
	}

	str := make([]byte, strlen)
	_, _ = buf.Read(str)
	return string(str), offset, nil
}

func readCert(buf *bytes.Reader) (*Cert, error) {
	var (
		offset int64
		err    error
		cert   = new(Cert)
	)

	cert.Alias, offset, err = readStr(buf, "certificate alias")
	if err != nil {
		return nil, err
	}

	cert.Timestamp, _, err = readTimestamp(buf)
	if err != nil {
		return nil, err
	}

	certType, _, err := readStr(buf, "certificate type")
	if certType != CertType {
		return nil, fmt.Errorf("unexpected certificate type at "+
			"position %d; found %q, expected %q",
			offset, certType, CertType)
	}

	elen, _, err := readUint32(buf, "encoded certificate length")
	if err != nil {
		return nil, err
	}

	if buf.Len() < int(elen) {
		return nil, fmt.Errorf("not enough data to read "+
			"certificate %q at position %d (length %d bytes)",
			cert.Alias, offset, elen)
	}

	cert.Raw = make([]byte, elen)
	_, _ = buf.Read(cert.Raw)

	cert.Cert, cert.CertErr = x509.ParseCertificate(cert.Raw)
	return cert, nil
}

func readKeypair(buf *bytes.Reader, opts *Options) (*Keypair, error) {
	var (
		offset   int64
		err      error
		certType string
		kp       = new(Keypair)
	)

	// retrive the key's alias, and use this to search for a password
	kp.Alias, offset, err = readStr(buf, "certificate alias")
	if err != nil {
		return nil, err
	}
	passwd, ok := opts.KeyPasswords[kp.Alias]
	if !ok {
		// no specific password for this alias, so use the file password
		passwd = opts.Password
	}

	kp.Timestamp, _, err = readTimestamp(buf)
	if err != nil {
		return nil, err
	}

	elen, _, err := readUint32(buf, "encrypted private key length")
	if err != nil {
		return nil, err
	}

	if buf.Len() < int(elen) {
		return nil, fmt.Errorf("not enough data to read "+
			"private key %q at position %d (length %d bytes)",
			kp.Alias, offset, elen)
	}

	kp.EncryptedKey = make([]byte, elen)
	_, _ = buf.Read(kp.EncryptedKey)
	kp.RawKey, kp.PrivKeyErr = DecryptPKCS8(kp.EncryptedKey, passwd)
	if kp.PrivKeyErr == nil {
		// we should now have a PKCS#8 PrivateKeyInfo, which Go can
		// parse for us
		kp.PrivateKey, kp.PrivKeyErr = x509.ParsePKCS8PrivateKey(
			kp.RawKey)
	}

	ncerts, _, err := readUint32(buf, "length of certificate chain")
	if err != nil {
		return nil, err
	}

	for n := uint32(0); n < ncerts; n++ {
		certType, offset, err = readStr(buf, fmt.Sprintf(
			"certificate type (chain entry #%d for %q)",
			n+1, kp.Alias))
		if err != nil {
			return nil, err
		}
		if certType != CertType {
			return nil, fmt.Errorf("unexpected certificate type "+
				"%q (expected %q at position %d for chain "+
				"entry #%d for %q)",
				certType, CertType, offset, n+1, kp.Alias)
		}

		elen, _, err = readUint32(buf, fmt.Sprintf(
			"encoded certificate length (chain entry #%d for %q)",
			n+1, kp.Alias))
		if err != nil {
			return nil, err
		}

		if buf.Len() < int(elen) {
			return nil, fmt.Errorf("not enough data to read "+
				"certificate chain entry #%d for %q at "+
				"position %d (length %d bytes)",
				n+1, kp.Alias, offset, elen)
		}

		kpc := new(KeypairCert)
		kpc.Raw = make([]byte, elen)
		_, _ = buf.Read(kpc.Raw)
		kpc.Cert, kpc.CertErr = x509.ParseCertificate(kpc.Raw)

		kp.CertChain = append(kp.CertChain, kpc)
	}

	return kp, nil
}
