package hash

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMD5URLEncoder_Encode(t *testing.T) {
	e := NewMD5URLEncoder()

	url, err := e.Encode("https://google.com", primitive.NewObjectID(), 0, 8)

	require.NoError(t, err)
	require.NotNil(t, url)
}
