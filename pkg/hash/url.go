package hash

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrURLAliasLengthExceed = errors.New("cannot generate alias due to length")

type URLEncoder interface {
	Encode(url string, userId primitive.ObjectID, start int, length int) (string, error)
}

type MD5URLEncoder struct {
}

func NewMD5URLEncoder() *MD5URLEncoder {
	return &MD5URLEncoder{}
}

func (e *MD5URLEncoder) Encode(url string, userId primitive.ObjectID, start int, length int) (string, error) {
	hasher := md5.New()

	_, err := hasher.Write([]byte(url + userId.Hex()))

	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	encoded := base64.StdEncoding.EncodeToString([]byte(hash))

	if start+length > len(encoded) {
		return "", ErrURLAliasLengthExceed
	}

	return encoded[start : start+length], nil
}
