package hash

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

type URLEncoder interface {
	Encode(url string) (string, error)
}

type MD5Encoder struct {
}

func NewMD5Encoder() *MD5Encoder {
	return &MD5Encoder{}
}

func (e *MD5Encoder) Encode(url string) (string, error) {
	hasher := md5.New()

	_, err := hasher.Write([]byte(url))

	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	encoded := base64.StdEncoding.EncodeToString([]byte(hash))

	return encoded[:8], nil
}
