package sha3x

import (
	"encoding/base64"
	"golang.org/x/crypto/sha3"
)

func Sum256Base64(data []byte, encoder *base64.Encoding) string {
	h := sha3.Sum256(data)
	return encoder.EncodeToString(h[:])
}
