package vault_test

import (
	"testing"
	"../vault"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigurationSuccess(t *testing.T) {

	expectedConfig := vault.Config {
		UserName: "user",
		Password: "password",
		KubeToken: "myCrazyKubeToken",
	}

	expectedConfig.Vault = vault.Vault {
		Server: "http://localhost:8200",
		Role: "myrole",
		LoginPath: "kubernetes",
		AuthMethod: "kubernetes",
		SecretPath: "demokv/mytest",
		KvVersion: "v1",
	}

	configPath := "testfiles/TestConfiguration.yaml"
	config, err := vault.ReadConfig(configPath)

	assert.Equal(t, config.UserName, expectedConfig.UserName)
	assert.Equal(t, config.Password, expectedConfig.Password)
	assert.Equal(t, config.KubeToken, expectedConfig.KubeToken)
	assert.Equal(t, config.Vault, expectedConfig.Vault)
	assert.NoError(t, err)
	
}
func TestReadConfigurationFromFileFail(t *testing.T) {
	
	configPath := "notexisting.yaml"
	config, err := vault.ReadConfig(configPath)

	assert.Nil(t, config)
	assert.Error(t, err)
}