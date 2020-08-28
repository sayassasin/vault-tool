package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"./vault"
	"strings"
)

func main() {

	configFilePtr := flag.String("config", "", "name of the configfile")
	valuesPtr := flag.String("secrets", "", "file with values to write")
	writePtr := flag.Bool("write", false, "write values to secretPath from config. Requires secrets parameter to be set")

	flag.Parse()

	if *configFilePtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *writePtr == true && *valuesPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	configFile := *configFilePtr
	fmt.Println("using config from: " + configFile)

	config, err := vault.ReadConfig(configFile)
	if err != nil {
		log.Fatal("error reading config")
	}
	if strings.ToLower(config.Vault.AuthMethod) == "kubernetes" {
		vault.KubernetesLogin(config)
	} else if strings.ToLower(config.Vault.AuthMethod) == "userpass" {
		vault.UserpassLogin(config)
	} else {
		log.Fatal("unsupported login method")
	}

	
	if *writePtr == false {
		vaultData := vault.ReadSecretsFromVault(config)
		for key, value := range vaultData {
			fmt.Println(key+": ", value.(string))
		}
	} else {
		secrets, err := vault.ReadSecretsFromFile(*valuesPtr)
		if err != nil {
			log.Fatal("error reading file contianing secrets to write")
		}
		if strings.EqualFold(strings.TrimSpace(secrets.Mode), "UPDATE") { 
			existingData := vault.ReadSecretsFromVault(config)

			for key, value := range existingData {
				if _, ok := secrets.Values[key]; ok {
					continue
				} else {
					secrets.Values[key] = fmt.Sprint(value)
				}
			}
		} else if strings.EqualFold(strings.TrimSpace(secrets.Mode), "WRITE") {
			vault.WriteSecrets(config, secrets)
		}
	}
}
