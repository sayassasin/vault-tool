package vault_test

import (
	"testing"
	"io/ioutil"
	"bytes"
	"net/http"
	"../restclient"
	"../mocks"
	"../vault"
	"github.com/stretchr/testify/assert"
)

func init() {
	restclient.Client = &mocks.MockClient{}
}

func TestUserpassLoginSuccess(t *testing.T) {
	// build response JSON
	json := `{"request_id":"fd8200c8-74e8-4c71-e935-f6d4312efd7a","lease_id":"","renewable":false,"lease_duration":0,"data":null,"wrap_info":null,"warnings":null,"auth":{"client_token":"s.ejSPQIGVAIgXTbS8HGWYCfgn","accessor":"GApEYTqoO59wiKTEKNcDeeKG","policies":["default","my.new.policy"],"token_policies":["default"],"identity_policies":["my.new.policy"],"metadata":{"username":"user"},"lease_duration":2764800,"renewable":true,"entity_id":"f519aea6-e873-75c7-4b8c-a467e0d50617","token_type":"service","orphan":true}}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	config := vault.Config{UserName: "user", Password: "password"}
	config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1"}
	vault.UserpassLogin(&config)
	assert.Equal(t, "s.ejSPQIGVAIgXTbS8HGWYCfgn", config.Vault.VaultToken)
}

func TestUserpassLoginFail(t *testing.T) {
		// build response JSON
		json := `{"errors":["invalid username or password"]}`
		// create a new reader with that JSON
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       r,
			}, nil
		}
		config := vault.Config{UserName: "user", Password: "password"}
		config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1"}
		err := vault.UserpassLogin(&config)
		
		assert.Equal(t, "", config.Vault.VaultToken)
		assert.NotNil(t, err)
}

func TestKubernetesLoginSuccess(t *testing.T) {
		// build response JSON
		json := `{"request_id":"fd8200c8-74e8-4c71-e935-f6d4312efd7a","lease_id":"","renewable":false,"lease_duration":0,"data":null,"wrap_info":null,"warnings":null,"auth":{"client_token":"s.ejSPQIGVAIgXTbS8HGWYCfgn","accessor":"GApEYTqoO59wiKTEKNcDeeKG","policies":["default","my.new.policy"],"token_policies":["default"],"identity_policies":["my.new.policy"],"metadata":{"username":"user"},"lease_duration":2764800,"renewable":true,"entity_id":"f519aea6-e873-75c7-4b8c-a467e0d50617","token_type":"service","orphan":true}}`
		// create a new reader with that JSON
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		}
		config := vault.Config{KubeToken: "token"}
		config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1", Role: "testrole"}
		vault.KubernetesLogin(&config)
		assert.Equal(t, "s.ejSPQIGVAIgXTbS8HGWYCfgn", config.Vault.VaultToken)
}

func TestKubernetesLoginFail(t *testing.T) {
	// build response JSON
	json := `{"errors":["missing jwt"]}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, nil
	}
	config := vault.Config{KubeToken: ""}
	config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1", Role: "testrole"}
	err := vault.KubernetesLogin(&config)
	assert.Equal(t, "", config.Vault.VaultToken)
	assert.NotNil(t, err)
}

func TestReadSecretsFromVault(t *testing.T) {
	json := `{"request_id":"c7e16dea-73f1-7af7-58b7-3652c41c1c75","lease_id":"","renewable":false,"lease_duration":2764800,"data":{"test":"demo"},"wrap_info":null,"warnings":null,"auth":null}`

	reader := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       reader,
		}, nil
	}

	config := vault.Config{UserName: "dummy"}
	config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1", VaultToken: "abcdefgh"}

	var vaultData = vault.ReadSecretsFromVault(&config)
	assert.Equal(t, "demo", vaultData["test"])

}

func TestReadSecretsFromFileSuccess(t *testing.T) {
	values := map[string]string{
		"demo": "real",
		"cat": "dog",
		"hello": "world",
	} 
	expectedSecrets := vault.SecretsToWrite{
		Mode: "UPDATE",
		Values: values,
	}

	valuePath := "testfiles/TestValuesToWrite.yaml"
	secrets, err := vault.ReadSecretsFromFile(valuePath)

	assert.Equal(t, secrets.Mode, expectedSecrets.Mode)
	assert.Equal(t, secrets.Values, expectedSecrets.Values)
	assert.NoError(t, err)
}

func TestReadSecretsFromFileFail(t *testing.T) {

	valuePath := "notexisting.yaml"
	secrets, err := vault.ReadSecretsFromFile(valuePath)

	assert.Nil(t, secrets)
	assert.Error(t, err)
}

func TestWriteSecretsSuccess(t *testing.T) {
	values := map[string]string{
		"demo": "real",
		"cat": "dog",
		"hello": "world",
	} 
	secrets := vault.SecretsToWrite{
		Mode: "UPDATE",
		Values: values,
	}
	config := vault.Config{UserName: "user", Password: "password"}
	config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1", VaultToken: "abc"}
	
	reader := ioutil.NopCloser(bytes.NewReader([]byte("")))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       reader,
		}, nil
	}
	
	err := vault.WriteSecrets(&config, &secrets)
	assert.NoError(t, err)

}

func TestWriteSecretsFail(t *testing.T) {
	values := map[string]string{
		"demo": "real",
		"cat": "dog",
		"hello": "world",
	} 
	secrets := vault.SecretsToWrite{
		Mode: "UPDATE",
		Values: values,
	}
	config := vault.Config{UserName: "user", Password: "password"}
	config.Vault = vault.Vault{Server: "test", LoginPath: "path", SecretPath: "secretpath", KvVersion: "v1", VaultToken: "abc"}
	
	json := `{"errors":["permission denied"]}`

	reader := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 400,
			Body:       reader,
		}, nil
	}
	
	err := vault.WriteSecrets(&config, &secrets)
	assert.Error(t, err)

}