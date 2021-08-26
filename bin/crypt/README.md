# crypt

## Install

```bash
go install github.com/bketelsen/crypt/bin/crypt@latest
```

## Backends

crypt supports etcd and consul as backends via the `-backend` flag.

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

```bash
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

```bash
crypt get -secret-keyring secring.gpg /app/config
```

### Support for unencrypted values

```bash
crypt set -plaintext ...
crypt get -plaintext ...
```

Crypt now has support for getting and setting plain unencrypted values, as
a convenience.  This was added to the backend libraries so it could be exposed
in spf13/viper. Use the -plaintext flag to get or set a value without encryption. 
