package cryptstate

import (
	"testing"
)

func BlockCompare(a []byte, b []byte) (match bool) {
	if len(a) != len(b) {
		return
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return
		}
	}

	match = true
	return
}

func TestTimes2(t *testing.T) {
	msg := [AESBlockSize]byte{
		0x80, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	}
	expected := [AESBlockSize]byte{
		0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7b,
	}

	times2(msg[0:])
	if BlockCompare(msg[0:], expected[0:]) == false {
		t.Errorf("times2 produces invalid output: %v, expected: %v", msg, expected)
	}
}

func TestTimes3(t *testing.T) {
	msg := [AESBlockSize]byte{
		0x80, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	}
	expected := [AESBlockSize]byte{
		0x81, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x85,
	}

	times3(msg[0:])
	if BlockCompare(msg[0:], expected[0:]) == false {
		t.Errorf("times3 produces invalid output: %v, expected: %v", msg, expected)
	}
}

func TestZeros(t *testing.T) {
	var msg [AESBlockSize]byte
	zeros(msg[0:])
	for i := 0; i < len(msg); i++ {
		if msg[i] != 0 {
			t.Errorf("zeros does not zero slice.")
		}
	}
}

func TestXor(t *testing.T) {
	msg := [AESBlockSize]byte{
		0x80, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	}
	var out [AESBlockSize]byte
	xor(out[0:], msg[0:], msg[0:])
	for i := 0; i < len(out); i++ {
		if out[i] != 0 {
			t.Errorf("XOR broken")
		}
	}
}

func TestEncrypt(t *testing.T) {
	msg := [15]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}
	key := [AESBlockSize]byte{
		0x96, 0x8b, 0x1b, 0x0c, 0x53, 0x1e, 0x1f, 0x80, 0xa6, 0x1d, 0xcb, 0x27, 0x94, 0x09, 0x6f, 0x32,
	}
	eiv := [AESBlockSize]byte{
		0x1e, 0x2a, 0x9b, 0xd0, 0x2d, 0xa6, 0x8e, 0x46, 0x26, 0x85, 0x83, 0xe9, 0x14, 0x2a, 0xff, 0x2a,
	}
	div := [AESBlockSize]byte{
		0x73, 0x99, 0x9d, 0xa2, 0x03, 0x70, 0x00, 0x96, 0xef, 0x55, 0x06, 0x7a, 0x8b, 0xbe, 0x00, 0x07,
	}
	expected := [19]byte{
		0x1f, 0xfc, 0xdd, 0xb4, 0x68, 0x13, 0x68, 0xb7, 0x92, 0x67, 0xca, 0x2d, 0xba, 0xb7, 0x0d, 0x44, 0xdf, 0x32, 0xd4,
	}
	expected_eiv := [AESBlockSize]byte{
		0x1f, 0x2a, 0x9b, 0xd0, 0x2d, 0xa6, 0x8e, 0x46, 0x26, 0x85, 0x83, 0xe9, 0x14, 0x2a, 0xff, 0x2a,
	}

	cs, err := New()
	if err != nil {
		t.Errorf("%v", err)
	}

	out := make([]byte, 19)
	cs.SetKey(key[0:], eiv[0:], div[0:])
	cs.Encrypt(out[0:], msg[0:])

	if BlockCompare(out[0:], expected[0:]) == false {
		t.Errorf("Mismatch in output")
	}

	if BlockCompare(cs.EncryptIV[0:], expected_eiv[0:]) == false {
		t.Errorf("EIV mismatch")
	}
}

func TestDecrypt(t *testing.T) {
	key := [AESBlockSize]byte{
		0x96, 0x8b, 0x1b, 0x0c, 0x53, 0x1e, 0x1f, 0x80, 0xa6, 0x1d, 0xcb, 0x27, 0x94, 0x09, 0x6f, 0x32,
	}
	eiv := [AESBlockSize]byte{
		0x1e, 0x2a, 0x9b, 0xd0, 0x2d, 0xa6, 0x8e, 0x46, 0x26, 0x85, 0x83, 0xe9, 0x14, 0x2a, 0xff, 0x2a,
	}
	div := [AESBlockSize]byte{
		0x73, 0x99, 0x9d, 0xa2, 0x03, 0x70, 0x00, 0x96, 0xef, 0x55, 0x06, 0x7a, 0x8b, 0xbe, 0x00, 0x07,
	}
	crypted := [19]byte{
		0x1f, 0xfc, 0xdd, 0xb4, 0x68, 0x13, 0x68, 0xb7, 0x92, 0x67, 0xca, 0x2d, 0xba, 0xb7, 0x0d, 0x44, 0xdf, 0x32, 0xd4,
	}
	expected := [15]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}
	post_div := [AESBlockSize]byte{
		0x1f, 0x2a, 0x9b, 0xd0, 0x2d, 0xa6, 0x8e, 0x46, 0x26, 0x85, 0x83, 0xe9, 0x14, 0x2a, 0xff, 0x2a,
	}

	cs, err := New()
	if err != nil {
		t.Errorf("%v", err)
	}

	out := make([]byte, 15)
	cs.SetKey(key[0:], div[0:], eiv[0:])
	cs.Decrypt(out[0:], crypted[0:])

	if BlockCompare(out[0:], expected[0:]) == false {
		t.Errorf("Mismatch in output")
	}

	if BlockCompare(cs.DecryptIV[0:], post_div[0:]) == false {
		t.Errorf("Mismatch in DIV")
	}
}
