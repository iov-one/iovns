package utils

import (
	"encoding/base64"
)

// Base64Encode
func Base64Encode(str string) []byte {
	encodedLen := base64.RawStdEncoding.EncodedLen(len(str))
	encoded := make([]byte, encodedLen)
	base64.RawStdEncoding.Encode(encoded, []byte(str))
	return encoded
}

// Base64Decode decodes
func Base64Decode(key []byte) (string, error) {
	decodedLen := base64.RawStdEncoding.DecodedLen(len(key))
	decoded := make([]byte, decodedLen)
	_, err := base64.RawStdEncoding.Decode(decoded, key)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
