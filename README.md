# crypt


## Usage

### Encrypted and store a value

```
crypt set -keyring pubring.gpg /app/config config.json 
```

### Retrieve and decrypted value

```
crypt get -secret-keyring secring.gpg /app/config
```
