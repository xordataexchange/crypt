package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var flagset = flag.NewFlagSet("crypt", flag.ExitOnError)

var (
	data          string
	backend       string
	key           string
	keyring       string
	endpoint      string
	secretKeyring string
)

func init() {
	flagset.StringVar(&backend, "backend", "etcd", "backend provider")
	flagset.StringVar(&endpoint, "endpoint", "http://127.0.0.1:4001", "backend url")
}

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		help()
	}
	cmd := os.Args[1]
	switch cmd {
	case "set":
		setCmd(flagset)
	case "get":
		getCmd(flagset)
	default:
		help()
	}
}

func help() {
	fmt.Fprintf(os.Stderr, "usage: %s COMMAND [arg...]", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n\n")
	fmt.Fprintf(os.Stderr, "commands:\n")
	fmt.Fprintf(os.Stderr, "   get  retrieve the value of a key\n")
	fmt.Fprintf(os.Stderr, "   set  set the value of a key\n")
	os.Exit(1)
}
