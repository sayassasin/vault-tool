package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"../restclient"
)

type SecretsToWrite struct {
	Mode   string            `yaml:"mode"`
	Values map[string]string `yaml:"values"`
}

func UserpassLogin(config *Config) (error){
	vaultLoginPath := strings.Join([]string{config.Vault.Server, "v1", "auth", config.Vault.LoginPath, "login", config.UserName}, "/")
	var requestBody = map[string]interface{} {
		"password": config.Password,
	}

	response, err := performLogin(vaultLoginPath, requestBody)
	if err != nil {
		return err
	}
	unmarshalClientToken(response, config)
	return nil
}

func KubernetesLogin(config *Config) error {
	vaultLoginPath := strings.Join([]string{config.Vault.Server, "v1", "auth", config.Vault.LoginPath, "login"}, "/")
	var requestBody = map[string]interface{} {
		"jwt":  config.KubeToken,
		"role": config.Vault.Role,
	}
	//fmt.Print(vaultLoginPath)
	response, err := performLogin(vaultLoginPath, requestBody)
	if err != nil {
		return err
	}
	unmarshalClientToken(response, config)
	return nil
}

func ReadSecretsFromVault(config *Config) map[string]interface{} {
	vaultSecretPath := strings.Join([]string{config.Vault.Server, config.Vault.KvVersion, config.Vault.SecretPath}, "/")
	fmt.Println("Querying secret: " + vaultSecretPath)
	request, _ := http.NewRequest("GET", vaultSecretPath, nil)
	request.Header.Set("X-Vault-Token", config.Vault.VaultToken)
	response, err := restclient.Client.Do(request)
	if err != nil {
		log.Fatal("error retrieving secrets")
	}
	data, err := ioutil.ReadAll(response.Body)

	response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("Error! Statuscode %d retrieved", response.StatusCode)
	}
	var resp map[string]interface{}
	json.Unmarshal(data, &resp)
	vaultData := resp["data"].(map[string]interface{})

	for key, value := range vaultData {
		fmt.Println(key+": ", value.(string))
	}

	return vaultData
}

func ReadSecretsFromFile(valuesPath string) (*SecretsToWrite, error) {
	secrets := &SecretsToWrite{}

	file, err := os.Open(valuesPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

func WriteSecrets(config *Config, secrets *SecretsToWrite) error{
	var vaultSecretPath = strings.Join([]string{config.Vault.Server, config.Vault.KvVersion, config.Vault.SecretPath}, "/")
	fmt.Println("Writing secrets to: " + vaultSecretPath)

	requestBody, err := json.Marshal(secrets.Values)

	//fmt.Println(string(requestBody))

	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", vaultSecretPath, bytes.NewBuffer(requestBody))
	request.Header.Set("X-Vault-Token", config.Vault.VaultToken)
	response, err := restclient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("error retrieving secrets")
	}

	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("Statuscode %d retrieved", response.StatusCode)
	}
	
	log.Printf("Secrets written! HTTP %d returned.", response.StatusCode)
	return nil
}

func unmarshalClientToken(data []byte, config *Config) {
	var resp map[string]interface{}
	json.Unmarshal(data, &resp)
	auth := resp["auth"].(map[string]interface{})
	config.Vault.VaultToken = auth["client_token"].(string)
	if config.Vault.VaultToken == "" {
		log.Fatalf("Error parsing vaultToken: %s", resp)
	}
}

func performLogin(loginPath string, requestBody map[string]interface{}) ([]byte, error) {
	response, err := restclient.Post(loginPath, requestBody, nil)
	if err != nil {
		fmt.Println(err)
		log.Fatal("error performing request")
	}
	data, err := ioutil.ReadAll(response.Body)

	response.Body.Close()
	if response.StatusCode != 200 {
		log.Printf("Error! Statuscode %d retrieved", response.StatusCode)
		return nil, fmt.Errorf("Error! Statuscode %d retrieved", response.StatusCode)
	} else if strings.Contains(string(data), "errors") {
		return nil, fmt.Errorf("Error during login! %s", string(data))
	} else {
		return data, nil
	}
}

