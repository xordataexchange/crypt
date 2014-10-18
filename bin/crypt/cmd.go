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
	kr, err := os.Open(secretKeyring)
	if err != nil {
		log.Fatal(err)
	}
	defer kr.Close()
	data, err := backendStore.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	value, err := secconf.Decode(data, kr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", value)
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
	kr, err := os.Open(keyring)
	if err != nil {
		log.Fatal(err)
	}
	defer kr.Close()
	secureValue, err := secconf.Encode(d, kr)
	if err != nil {
		log.Fatal(err)
	}
	if err := backendStore.Set(key, secureValue); err != nil {
		log.Fatal(err)
	}
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
