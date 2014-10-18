package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xordataexchange/crypt/backend/etcd"
	"github.com/xordataexchange/crypt/config"
	"github.com/xordataexchange/crypt/encoding/secconf"
)

func getCmd(flagset *flag.FlagSet) {
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s get [args...] key\n", os.Args[0])
		flagset.PrintDefaults()
	}
	flagset.StringVar(&secretKeyring, "secret-keyring", ".secring.gpg", "path to secret keyring")
	flagset.Parse(os.Args[2:])
	key := flagset.Arg(0)
	if key == "" {
		log.Fatal("a key is required")
	}
	skr, err := os.Open(secretKeyring)
	if err != nil {
		log.Fatal(err)
	}
	defer skr.Close()
	machines := []string{endpoint}
	cm := config.NewEtcdConfigManager(machines, skr)
	value, err := cm.Get(key)
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
	flagset.StringVar(&keyring, "keyring", ".pubring.gpg", "path to public keyring")
	flagset.Parse(os.Args[2:])
	key := flagset.Arg(0)
	if key == "" {
		log.Fatal("a key is required")
	}
	data := flagset.Arg(1)
	if data == "" {
		log.Fatal("a data file is required")
	}
	backend := etcd.New([]string{endpoint})
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
	if err := backend.Set(key, secureValue); err != nil {
		log.Fatal(err)
	}
}
