package kubeconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	name := "issue-113"
	endpoint := "remote"
	port := "32217"
	assert := require.New(t)
	contents, err := os.ReadFile("testdata/admin.kubeconfig")
	assert.NoError(err)
	config, err := configureKubeConfig(contents, name, endpoint, port)
	assert.NoError(err)
	assert.NotEmpty(config.Clusters[name], "expected to find cluster with changed named")
	assert.True(config.Clusters[name].InsecureSkipTLSVerify, "expected to find insecure access setup")
	assert.Nil(config.Clusters[name].CertificateAuthorityData, "expected to not find any certificate-authority-data")
}
