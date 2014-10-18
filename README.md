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

### As a library

```
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xordataexchange/crypt/config"
)

var (
	key           = "/app/config"
	secretKeyring = ".secring.gpg"
)

func main() {
	skr, err := os.Open(secretKeyring)
	if err != nil {
		log.Fatal(err)
	}
	defer skr.Close()
	machines := []string{"http://127.0.0.1:4001"}
	cm := config.NewEtcdConfigManager(machines, skr)
	value, err := cm.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", value)
}
```
