package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ImageLifeCycle(t *testing.T) {
	assert := require.New(t)
	client, err := NewClient(context.TODO())
	assert.NoError(err)
	err = client.CreateImage("dev", "testdata/supportbundle_f159fbe2-dae7-4606-b81c-f54e1a562c99_2024-11-18T04-34-27Z.zip", "rancher/support-bundle-kit:master-head")
	assert.NoError(err)
	images, err := client.FindImages("dev")
	assert.NoError(err)
	assert.NotZero(images, 0, "expected to atleast find one image")
	err = client.RemoveImages("dev")
	assert.NoError(err)
}
