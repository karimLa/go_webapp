package lib

import (
	"crypto/rand"
	"encoding/base64"
)

const RememeberTokenBytes = 32

// Bytes will help us generate n random bytes, or will
// return an error if there was one. This uses the crypto/rand
// package so it is safe to use with things like rememeber tokens,
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Base64FromBytes will generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
// of that byte slice
func Base64FromBytes(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// RememeberToken is a helper function designed to generate
// remember tokens of a predetermined byte size.
func RememeberToken() (string, error) {
	return Base64FromBytes(RememeberTokenBytes)
}
