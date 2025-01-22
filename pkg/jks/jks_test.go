package jks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
	"unicode/utf16"
)

// TestComputeDigest is a regression test for the digest function.
func TestComputeDigest(t *testing.T) {
	t.Run("empty", testComputeDigest("", "",
		"569D05A766C473698C0B58EBAEAE0A25EB10BACC"))
	t.Run("regr", testComputeDigest("input data", "password",
		"74DDD13B68919674D4409A19AB284019A1DA57C8"))
}

func testComputeDigest(in, passwd, expHex string) func(*testing.T) {
	return func(t *testing.T) {
		exp, err := hex.DecodeString(expHex)
		if err != nil {
			t.Fatalf("error decoding expHex: %v", err)
		}
		out := ComputeDigest([]byte(in), passwd)
		if !bytes.Equal(out, exp) {
			t.Errorf("output sequence (len %d) ≠ expected",
				len(out))
			t.Errorf("out %X", out)
		}
	}
}

// TestPasswordUTF16 checks that our UTF-16 encoding routine works as expected.
// The test cases incorporate empty strings and Unicode strings with characters
// outside the BMP (basic multilingual plane), i.e. ones that need encoding as
// UTF-16 surrogate pairs.
func TestPasswordUTF16(t *testing.T) {
	t.Run("empty", testPasswordBytes("", nil))
	t.Run("ascii-1", testPasswordBytes("ascii",
		[]byte{0, 'a', 0, 's', 0, 'c', 0, 'i', 0, 'i'}))
	t.Run("ascii-2", testPasswordUTF16("ascii"))
	t.Run("utf8", testPasswordUTF16("a≤b"))
	t.Run("surrogate", testPasswordUTF16("z1\U00016000\u2340•—@.µ"))
}

func testPasswordBytes(in string, exp []byte) func(*testing.T) {
	return func(t *testing.T) {
		out := PasswordUTF16(in)
		if !bytes.Equal(out, exp) {
			t.Errorf("output sequence ‘%X’ ≠ expected ‘%X’",
				out, exp)
		}
	}
}

func testPasswordUTF16(in string) func(*testing.T) {
	return func(t *testing.T) {
		out := PasswordUTF16(in)
		expStr := utf16.Encode([]rune(in))
		exp := make([]byte, len(expStr)*2)
		for i, v := range expStr {
			binary.BigEndian.PutUint16(exp[i*2:], v)
		}
		if !bytes.Equal(out, exp) {
			t.Errorf("output sequence ‘%X’ ≠ expected ‘%X’",
				out, exp)
		}
	}
}
