package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"math/big"
)

func GenerateString(length int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func HashString(str, salt string) (string, error) {
	hasher := sha512.New()

	if _, err := hasher.Write([]byte(str + salt)); err != nil {
		return "", err
	}

	hashed := hasher.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hashed)

	return encoded, nil
}
