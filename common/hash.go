package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"
)

func HashToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashMD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	return md5Hash
}

func HashUniqueID() string {
	currentTime := time.Now().Format("20060102150405")

	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := make([]byte, 10)

	for i := range randomStr {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		randomStr[i] = chars[index.Int64()]
	}
	return currentTime + string(randomStr)
}
