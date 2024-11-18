package docker

import (
	"archive/tar"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BuildContextTar(t *testing.T) {
	assert := require.New(t)
	buf, err := BuildContextTar("testdata/supportbundle_f159fbe2-dae7-4606-b81c-f54e1a562c99_2024-11-18T04-34-27Z.zip", "rancher/support-bundle-kit:master")
	assert.NoError(err)
	tr := tar.NewReader(buf)
	var dockerFileFound bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(err, "expected no error while parsing file")
		if hdr.Name == "Dockerfile" {
			dockerFileFound = true
			break
		}
	}
	assert.True(dockerFileFound, "expected to find dockerfile")
}
