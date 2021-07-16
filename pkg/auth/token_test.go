package auth

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var (
	token  = ""
	userId = primitive.NewObjectID().String()
)

func TestNewJWTManager_emptyKey(t *testing.T) {
	_, err := NewJWTManager("")

	require.Errorf(t, err, "empty signing key")
}

func TestJWTManager_Issue(t *testing.T) {
	m, err := NewJWTManager("key")

	require.NoError(t, err)

	token, err = m.Issue(userId, time.Duration(1)*time.Hour)

	require.NoError(t, err)
	require.NotNil(t, token)
}

func TestJWTManager_Decode(t *testing.T) {
	m, err := NewJWTManager("key")

	require.NoError(t, err)

	id, err := m.Decode(token)

	require.NoError(t, err)
	require.Equal(t, userId, id)
}

func TestJWTManager_DecodeErr(t *testing.T) {
	m, err := NewJWTManager("key")

	require.NoError(t, err)

	_, err = m.Decode("qwe")

	require.Error(t, err)
}
