package bytes32_test

import (
	"testing"

	"github.com/Abso1ut3Zer0/Dexalot-go/src/bytes32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDecodeToString(t *testing.T) {
	str := "HELLO_WORLD"
	bytes := []byte(str)
	var b [32]byte
	copy(b[:], bytes)
	decoded := bytes32.DecodeToString(b)
	assert.Equal(t, str, decoded)

	str = "THIS_MESSAGE_IS_LONGER_THAN_32_BYTES"
	bytes = []byte(str)
	copy(b[:], bytes)
	decoded = bytes32.DecodeToString(b)
	assert.NotEqual(t, str, decoded)
	assert.Equal(t, len(decoded), 32)
	assert.Equal(t, decoded, "THIS_MESSAGE_IS_LONGER_THAN_32_B")
}

func TestDecodeHexToString(t *testing.T) {
	str := "HELLO_WORLD"
	bytes := []byte(str)
	var b1 [32]byte
	copy(b1[:], bytes)
	decoded, err := bytes32.DecodeHexToString(b1)
	assert.Nil(t, err)
	if decoded != str {
		t.Errorf("Expected %s, got %s", str, decoded)
	}

	str = "0x415641582f555344430000000000000000000000000000000000000000000000"
	hexed := common.HexToHash(str)
	bytes = hexed.Bytes()
	var b2 [32]byte
	copy(b2[:], bytes)
	decoded, err = bytes32.DecodeHexToString(b2)
	assert.Nil(t, err)
	assert.Equal(t, "AVAX/USDC", decoded)
}
