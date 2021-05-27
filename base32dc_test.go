package base32dc

import (
	"math/rand"
	"testing"
)

func Test_VerifyCheckSum(t *testing.T) {
	var srcBytes1 [1]byte
	var enc string
	srcBytes1[0] = 0
	enc = EncodeWithCheckSum(srcBytes1[:])
	if !VerifyCheckSum(enc) {
		t.Fail()
	}
	var srcBytes5 [5]byte
	for i := range srcBytes5 {
		srcBytes5[1] = byte(150 + i)
	}
	enc = EncodeWithCheckSum(srcBytes5[:])
	if !VerifyCheckSum(enc) {
		t.Fail()
	}
	enc = "_" + enc[1:]
	if VerifyCheckSum(enc) {
		t.Fail()
	}
	if enc[0] == 'a' {
		enc = "b" + enc[1:]
	} else {
		enc = "a" + enc[1:]
	}
	if VerifyCheckSum(enc) {
		t.Fail()
	}
}
func Test_Base32(t *testing.T) {
	var srcBytes [16]byte
	var enc string
	var dec [16]byte

	src := srcBytes[:]

	for i := range src {
		src[i] = byte(255 - i)
	}
	enc = Encode(src)
	if enc != "zqzvfyfztfyhzvvyn7x7fs7yg7" {
		t.Fail()
	}
	enc = EncodeWithCheckSum(src)
	if enc != "zqzvfyfztfyhzvvyn7x7fs7yg7~" {
		t.Fail()
	}
	if !VerifyCheckSum(enc) {
		t.Fail()
	}

	for i := 0; i < 1000; i++ {
		rand.Read(src)
		enc = Encode(src)
		if len(enc) != 26 {
			t.Fail()
		}
		if !Decode(enc, dec[:]) {
			t.Fail()
		}
		if srcBytes != dec {
			t.Fail()
		}
		enc = EncodeWithCheckSum(src)
		if len(enc) != 27 {
			t.Fail()
		}
		if !DecodeWithCheckSum(enc, dec[:]) {
			t.Fail()
		}
		if srcBytes != dec {
			t.Fail()
		}
	}
}
