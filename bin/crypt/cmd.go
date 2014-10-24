package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xordataexchange/crypt/backend"
	"github.com/xordataexchange/crypt/backend/consul"
	"github.com/xordataexchange/crypt/backend/etcd"
	"github.com/xordataexchange/crypt/encoding/secconf"
	"github.com/xordataexchange/crypt/encoding/standard"
)

func getCmd(flagset *flag.FlagSet) {
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s get [args...] key\n", os.Args[0])
		flagset.PrintDefaults()
	}
	flagset.StringVar(&secretKeyring, "secret-keyring", ".secring.gpg", "path to armored secret keyring")
	flagset.Parse(os.Args[2:])
	key := flagset.Arg(0)
	if key == "" {
		flagset.Usage()
		os.Exit(1)
	}
	backendStore, err := getBackendStore(backendName, endpoint)
	if err != nil {
		log.Fatal(err)
	}
	if plaintext {
		value, err := getPlain(key, backendStore)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", value)
		return
	}
	value, err := getEncrypted(key, secretKeyring, backendStore)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", value)
}

func getEncrypted(key, keyring string, store backend.Store) ([]byte, error) {
	var value []byte
	kr, err := os.Open(secretKeyring)
	if err != nil {
		return value, err
	}
	defer kr.Close()
	data, err := store.Get(key)
	if err != nil {
		return value, err
	}
	value, err = secconf.Decode(data, kr)
	if err != nil {
		return value, err
	}
	return value, err

}

func getPlain(key string, store backend.Store) ([]byte, error) {
	var value []byte
	data, err := store.Get(key)
	if err != nil {
		return value, err
	}
	value, err = standard.Decode(data)
	if err != nil {
		return value, err
	}
	return value, err
}

func setCmd(flagset *flag.FlagSet) {
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s set [args...] key file\n", os.Args[0])
		flagset.PrintDefaults()
	}
	flagset.StringVar(&keyring, "keyring", ".pubring.gpg", "path to armored public keyring")
	flagset.Parse(os.Args[2:])
	key := flagset.Arg(0)
	if key == "" {
		flagset.Usage()
		os.Exit(1)
	}
	data := flagset.Arg(1)
	if data == "" {
		flagset.Usage()
		os.Exit(1)
	}
	backendStore, err := getBackendStore(backendName, endpoint)
	if err != nil {
		log.Fatal(err)
	}
	d, err := ioutil.ReadFile(data)
	if err != nil {
		log.Fatal(err)
	}

	if plaintext {
		err := setPlain(key, backendStore, d)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	err = setEncrypted(key, keyring, d, backendStore)
	if err != nil {
		log.Fatal(err)
	}
	return

}
func setPlain(key string, store backend.Store, d []byte) error {
	value, err := standard.Encode(d)
	if err != nil {
		return err
	}
	err = store.Set(key, value)
	return err

}

func setEncrypted(key, keyring string, d []byte, store backend.Store) error {
	kr, err := os.Open(keyring)
	if err != nil {
		return err
	}
	defer kr.Close()
	secureValue, err := secconf.Encode(d, kr)
	if err != nil {
		return err
	}
	err = store.Set(key, secureValue)
	return err
}

func getBackendStore(provider string, endpoint string) (backend.Store, error) {
	if endpoint == "" {
		switch provider {
		case "consul":
			endpoint = "127.0.0.1:8500"
		case "etcd":
			endpoint = "http://127.0.0.1:4001"
		}
	}
	machines := []string{endpoint}
	switch provider {
	case "etcd":
		return etcd.New(machines)
	case "consul":
		return consul.New(machines)
	default:
		return nil, errors.New("invalid backend " + provider)
	}
}
