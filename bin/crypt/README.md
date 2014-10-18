# crypt

## Install

```
go install github.com/xordataexchange/crypt/bin/crypt
```

## Usage

```
usage: crypt COMMAND [arg...]

commands:
   get  retrieve the value of a key
   set  set the value of a key
```

### Encrypted and set a value

```
usage: crypt set [args...] key file
  -backend="etcd": backend provider
  -endpoint="": backend url
  -keyring=".pubring.gpg": path to armored public keyring
```

Example:

```
crypt set -keyring pubring.gpg /app/config config.json 
```

### Retrieve and decrypted a value

```
usage: crypt get [args...] key
  -backend="etcd": backend provider
  -endpoint="": backend url
  -secret-keyring=".secring.gpg": path to armored secret keyring
```

Example:

```
crypt get -secret-keyring secring.gpg /app/config
```
