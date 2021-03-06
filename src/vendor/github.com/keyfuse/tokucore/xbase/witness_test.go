// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func witnessScriptPubkey(version byte, program []byte) []byte {
	if version != 0 {
		version += 0x50
	}
	return append(append([]byte{version}, byte(len(program))), program...)
}

func TestWitnessAddress(t *testing.T) {
	tests := []struct {
		address        string
		witnessProgram []byte
	}{
		{"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
			[]byte{
				0x00, 0x14, 0x75, 0x1e, 0x76, 0xe8, 0x19, 0x91, 0x96, 0xd4, 0x54,
				0x94, 0x1c, 0x45, 0xd1, 0xb3, 0xa3, 0x23, 0xf1, 0x43, 0x3b, 0xd6,
			},
		},
		{"tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7",
			[]byte{
				0x00, 0x20, 0x18, 0x63, 0x14, 0x3c, 0x14, 0xc5, 0x16, 0x68, 0x04,
				0xbd, 0x19, 0x20, 0x33, 0x56, 0xda, 0x13, 0x6c, 0x98, 0x56, 0x78,
				0xcd, 0x4d, 0x27, 0xa1, 0xb8, 0xc6, 0x32, 0x96, 0x04, 0x90, 0x32,
				0x62,
			},
		},
		{"bc1pw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7k7grplx",
			[]byte{
				0x51, 0x28, 0x75, 0x1e, 0x76, 0xe8, 0x19, 0x91, 0x96, 0xd4, 0x54,
				0x94, 0x1c, 0x45, 0xd1, 0xb3, 0xa3, 0x23, 0xf1, 0x43, 0x3b, 0xd6,
				0x75, 0x1e, 0x76, 0xe8, 0x19, 0x91, 0x96, 0xd4, 0x54, 0x94, 0x1c,
				0x45, 0xd1, 0xb3, 0xa3, 0x23, 0xf1, 0x43, 0x3b, 0xd6,
			},
		},
		{"bc1sw50qa3jx3s",
			[]byte{
				0x60, 0x02, 0x75, 0x1e,
			},
		},
		{"bc1zw508d6qejxtdg4y5r3zarvaryvg6kdaj",
			[]byte{
				0x52, 0x10, 0x75, 0x1e, 0x76, 0xe8, 0x19, 0x91, 0x96, 0xd4, 0x54,
				0x94, 0x1c, 0x45, 0xd1, 0xb3, 0xa3, 0x23,
			},
		},
		{"tb1qqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesrxh6hy",
			[]byte{
				0x00, 0x20, 0x00, 0x00, 0x00, 0xc4, 0xa5, 0xca, 0xd4, 0x62, 0x21,
				0xb2, 0xa1, 0x87, 0x90, 0x5e, 0x52, 0x66, 0x36, 0x2b, 0x99, 0xd5,
				0xe9, 0x1c, 0x6c, 0xe2, 0x4d, 0x16, 0x5d, 0xab, 0x93, 0xe8, 0x64,
				0x33,
			},
		},
	}

	for _, test := range tests {
		hrp, version, program, err := WitnessDecode(test.address)
		assert.Nil(t, err)
		assert.Equal(t, test.witnessProgram, witnessScriptPubkey(version, program))

		addr, err := WitnessEncode(hrp, version, program)
		assert.Nil(t, err)
		assert.Equal(t, test.address, addr)
	}
}
