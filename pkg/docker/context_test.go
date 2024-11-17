package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetCurrentContext(t *testing.T) {
	ctx := context.TODO()
	err := GetCurrentContext(ctx)
	assert := require.New(t)
	assert.NoError(err)
}
