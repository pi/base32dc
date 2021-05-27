// Douglas Crockford Base32 encoding
package base32dc

import (
	"crypto/rand"
	"strings"
)

const (
	kValueSymbols = "0123456789abcdefghjkmnpqrstvwxyz"
	kCheckSymbols = "*~$=u"
	kSymbols      = kValueSymbols + kCheckSymbols
)

var (
	kEncTable [37]byte // bits -> symbol, contains values for checksym symbols
	kDecTable [256]int // symbol -> bits, -1 -> invalid symbol
)

func init() {
	for i := 0; i <= 36; i++ {
		kEncTable[i] = byte(kSymbols[i])
	}
	for i := 0; i <= 255; i++ {
		kDecTable[i] = -1
	}
	for i := 0; i <= 36; i++ {
		kDecTable[int(kSymbols[i])] = i
		if kSymbols[i] >= 'a' && kSymbols[i] <= 'z' {
			kDecTable[int(kSymbols[i])+32] = i
		}
	}
	kDecTable[int('i')] = 1
	kDecTable[int('I')] = 1
	kDecTable[int('l')] = 1
	kDecTable[int('L')] = 1
	kDecTable[int('o')] = 0
	kDecTable[int('O')] = 0
}

func encode(val []byte, withCheckSum bool) string {
	vl := len(val)
	if vl == 0 {
		return ""
	}
	sb := strings.Builder{}
	for bi := 0; bi < vl*8; bi += 5 {
		byi := bi >> 3
		v := uint(val[byi])
		if byi != vl-1 {
			v |= uint(val[byi+1]) << 8
		}
		sb.WriteByte(kEncTable[(v>>(bi&0x7))&0x1f])
	}
	if withCheckSum {
		sb.WriteByte(kEncTable[val[0]%37])
	}
	return sb.String()
}

func Encode(val []byte) string {
	return encode(val, false)
}

func EncodeWithCheckSum(val []byte) string {
	return encode(val, true)
}

func decode(src string, dest []byte, withCheckSum bool) (result bool) {
	result = false
	n := len(src)
	if withCheckSum {
		n--
	}
	if n <= 0 {
		return
	}
	bitbuf := uint(0)
	nbits := 0
	di := 0
	var i int
	for i = 0; i < n && di < len(dest); i++ {
		v := kDecTable[int(src[i])]
		if v == -1 || v > 31 {
			return
		}
		bitbuf |= uint(v) << nbits
		nbits += 5
		if nbits >= 8 {
			dest[di] = byte(bitbuf & 0xFF)
			di++
			nbits -= 8
			bitbuf >>= 8
		}
	}
	if nbits > 0 {
		if di == len(dest) {
			if bitbuf != 0 {
				return
			}
		} else {
			dest[di] = byte(bitbuf)
			di++
		}
	}
	if i != n || di != len(dest) {
		return
	}

	if withCheckSum {
		return int(dest[0])%37 == kDecTable[int(src[n])]
	} else {
		return true
	}
}

func Decode(src string, dest []byte) bool {
	return decode(src, dest, false)
}

func DecodeWithCheckSum(src string, dest []byte) bool {
	return decode(src, dest, true)
}

func newGUID(withCheckSum bool) string {
	var uuid [16]byte
	rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return encode(uuid[:], withCheckSum)
}

func NewGUID() string {
	return newGUID(false)
}

func NewGUIDWithCheckSum() string {
	return newGUID(true)
}

func VerifyCheckSum(src string) bool {
	if len(src) < 2 {
		return false
	}
	for i := range src {
		dv := kDecTable[int(src[i])]
		if dv < 0 {
			return false
		}
		if dv > 31 {
			if i < len(src)-1 || dv > 36 {
				return false
			}
		}
	}
	lo := kDecTable[int(src[0])]
	var hi int
	if len(src) == 2 {
		hi = 0
	} else {
		hi = kDecTable[int(src[1])]
	}
	return (((lo | (hi << 5)) & 0xFF) % 37) == kDecTable[src[len(src)-1]]
}
