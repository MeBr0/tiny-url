package hash

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLEncoder interface {
	Encode(url string, userId primitive.ObjectID) (string, error)
}

type MD5Encoder struct {
}

func NewMD5Encoder() *MD5Encoder {
	return &MD5Encoder{}
}

func (e *MD5Encoder) Encode(url string, userId primitive.ObjectID) (string, error) {
	hasher := md5.New()

	_, err := hasher.Write([]byte(url + userId.Hex()))

	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	encoded := base64.StdEncoding.EncodeToString([]byte(hash))

	return encoded[:8], nil
}
