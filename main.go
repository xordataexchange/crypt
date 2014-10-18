package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xordataexchange/crypt/encoding/secconf"
)

var (
	data          string
	backend       string
	key           string
	keyring       string
	secretKeyring string
)

func init() {
	flag.StringVar(&backend, "backend", "", "backend")
	flag.StringVar(&data, "data", "", "path to configuration file")
	flag.StringVar(&key, "key", "", "configuration key")
	flag.StringVar(&keyring, "keyring", ".pubring.gpg", "path to public keyring")
	flag.StringVar(&secretKeyring, "secret-keyring", ".secring.gpg", "path to secret keyring")
}

var v = `wcBMA4vc/EzjLNNYAQgABXwWihzAYHtI7rFebu11NjsOhwsGhzHBAq0tJwHQGi7gCaz92ZOOA1e/1/GD7ZqRgskKqJd1KnMHvedYVRJbx0AriTouXICMSSpN9Q93GsXY6r7iVaozCCIAkW5YvTDwJ1/wkG9TMduoiflKglbV9LMBRObOb566FvGSlwivOF1eNYxQDbqUpaUXaw9QqQ9P/lRBho3J8Pn6Eg11dTclR9yirrc0IrnL7rBNZbZwC73ysVd3Oi7ZV24hPrDwld1zkoqMeoZQ0VuEF7W0tWmkOCVhQFuZO3xGtmeOnoCLXPMepeOHGlJCXsEAsjxcJUz0+x2hltfUqE0ld39AcUEGRNLgAeSBGfxTcDTbE15JSuaKXB3E4aZA4KXgv+FTD+Bn4iOsOLPgdeNIUS1vqKLU+eCK4e434ALlnw8qd3utvTbyxSb9zgOsbgVBxOBcoVIGz6g28yIPK+rgA+Qe8IE42t75mW+7G9oHatiS4IbilNgmBODB4U754Bfiv4mnGuCe4+enJLoAQp7K4NfkQxZ38e0ZkluHixQ34YWTgOJaj+vJ4Y/HAA==`

func main() {
	flag.Parse()
	cmd := "get"
	switch cmd {
	case "set":
		config, err := ioutil.ReadFile(data)
		if err != nil {
			log.Fatal(err)
		}
		kr, err := os.Open(keyring)
		if err != nil {
			log.Fatal(err)
		}
		defer kr.Close()
		secureValue, err := secconf.Encode(config, kr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("secure value: \n%s\n", secureValue)
		return
	case "get":
		skr, err := os.Open(secretKeyring)
		if err != nil {
			log.Fatal(err)
		}
		defer skr.Close()
		secureValue := []byte(v)
		value, err := secconf.Decode(secureValue, skr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("value: \n%s\n", value)
	default:
		log.Fatal("unknown command: ", cmd)
	}
}
