package docker

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ContainerLifeCycle(t *testing.T) {
	assert := require.New(t)
	client, err := NewClient(context.TODO())
	assert.NoError(err)
	err = client.CreateImage("issue-113", "testdata/supportbundle_f159fbe2-dae7-4606-b81c-f54e1a562c99_2024-11-18T04-34-27Z.zip", "rancher/support-bundle-kit:master-head")
	assert.NoError(err)
	err = client.RunContainer("issue-113", "testdata/supportbundle_f159fbe2-dae7-4606-b81c-f54e1a562c99_2024-11-18T04-34-27Z.zip")
	assert.NoError(err)
	contents, err := client.ReadFile("issue-7007", simKubeConfigPath)
	assert.NoError(err)
	file, err := os.CreateTemp(os.TempDir(), "kubeconfig")
	assert.NoError(err)
	err = os.WriteFile(file.Name(), contents, 0600)
	assert.NoError(err)
	assert.NotNil(contents, "expected content to not be nil")
	err = client.StopContainer("issue-113")
	assert.NoError(err)
	err = client.RemoveImages("issue-113")
	assert.NoError(err)
	assert.NoError(os.Remove(file.Name()), "expected no error while cleaning up temp file")
}
