package bytes32

import (
	"encoding/hex"
	"strings"
)

func DecodeToString(b [32]byte) string {
	str := string(b[:])
	str = strings.TrimRight(str, "\x00")
	return str
}

func DecodeHexToString(b [32]byte) (string, error) {
	hexStr := hex.EncodeToString(b[:])
	hexStr = strings.TrimLeft(hexStr, "0")
	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	str := string(hexBytes)
	str = strings.TrimRight(str, "\x00")
	return str, nil
}
