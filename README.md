# Vault-Tool

This tool supports the use of Hashicorp Vault. Vault CLI is great but having multiple secrets secured by different accounts can be tricky and Vault itself can be very complex for developers who doesn't care a lot about the infrastructure they use. At this point Vault-Tool joins the game. Just prepare different configuration files. For example one with userpass auth in to auth in your development environment and two with different Kubernetes service accounts for staging and production. The only difference is you just have to pass another file on your command line.

---
### Nerd Disclaimer ;)
This is my first tool written in Go-Lang :heart: (beside of the training stuff) and the code is sometimes a bit ugly (working on it!) For becoming better in Go I tried to avoid third party libraries like the official Vault library (and I only added yaml.V2 and testify for convenience ;) )

---
## Features

- Authenticate using userpass method
- Authenticate using kubernetes method
- Reading secrets from kv engine v1
- Writing secrets to kv engine v1
- Updating secrets in kv engine v1

---
## To Be Done

- Support for kv store v2
- More authentication methods
- Additional secret engines
- Refactoring code to write secrets

---
## How to use

### Authentication

To authenticate against Vault you have to create a YAML based configuration file. Below you can see examples for userpass and kubernetes authentification.
<br>

### Authenticating with userpass

```yaml
vault:
  server: "http://localhost:8200"
  secretPath: "demokv/mytest"
  loginPath: "userpass"
  authMethod: "userpass"
  kvVersion: "v1"
username: "sven"
password: "asdf"
```
<br>

### Authenticating with Kubernetes
You can get the kubetoken from your Kubernetes service account secret.

```sh
kubectl -n <your-namespace> get secret <your-secret-name> -o jsonpath="{.data.token}" | base64 --decode; echo;
```

```yaml
vault:
  server: "http://localhost:8200"
  role: "myrole"
  secretPath: "demokv/mytest"
  loginPath: "kubernetes"
  authMethod: "kubernetes"
  kvVersion: "v1"
kubeToken: "myCrazyKubeToken"
```
<br>

### Secret Management

Vault-Tool will perform http calls against the Vault api to read/write the secrets from the given secretPath.

### Reading secrets
To read a secret you only need a configuration.yaml an run Vault-Tool with the given config:
```sh
vaulttool -config myconfig.yaml
```
<br>

### Writing secrets

Writing secrets requires a second YAML file.

```yaml
mode: "WRITE"
values:
  demo: real
  cat: dog
  hello: world
```

This file has to be passed to vaulttool.
```sh
vaulttool -config myconfig.yaml -write true -secrets myValues.yaml
```
To avoid critical "oops.. moments" you also have to add the -write flag. else the secrets flag will be ignored.

<br>

### Updating secrets

Updating secrets is equal to writing secrets. You only have to change the mode in your secrets.yaml to <code>UPDATE</code>.

Vault-Tool will read the secrets from the given secret path an add new / update old secrets. It will not delete missing secrets. If you want to remove particular KVs you shoul you <code>WRITE</code>

```yaml
mode: "UPDATE"
values:
  demo: real
  cat: dog
  hello: world
```

## License
Take the code and do whatever you want! I am not responsible for any damage caused.
