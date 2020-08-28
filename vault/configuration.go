package vault

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Vault     Vault  `yaml:"vault""`
	KubeToken string `yaml:"kubeToken"`
	UserName  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type Vault struct {
	Server   	   string `yaml:"server""`
	Role 	       string `yaml:"role""`
	LoginPath      string `yaml:"loginPath"`
	SecretPath     string `yaml:"secretPath""`
	AuthMethod     string `yaml:"authMethod"`
	KvVersion      string `yaml:"kvVersion"`
	VaultToken     string
}

func ReadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	/*if config.Vault.LoginPath == "" {
		config.Vault.LoginPath = "userpass"
	}*/

	return config, nil
}
