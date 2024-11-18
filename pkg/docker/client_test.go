package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetClient(t *testing.T) {
	cli, err := GetClient()
	assert := require.New(t)
	assert.NoError(err)
	assert.NotNil(cli)
}
