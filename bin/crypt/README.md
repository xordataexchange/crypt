# crypt

## Usage

```
usage: crypt COMMAND [arg...]

commands:
   get  retrieve the value of a key
   set  set the value of a key
```

### Encrypted and set a value

```
crypt set -keyring pubring.gpg /app/config config.json 
```

### Retrieve and decrypted a value

```
crypt get -secret-keyring secring.gpg /app/config
```
