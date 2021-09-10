package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var flagset = flag.NewFlagSet("crypt", flag.ExitOnError)

var (
	Version    = "dev"
	Commit     string
	CommitDate string
	builtBy    string
)

var (
	data          string
	backendName   string
	key           string
	keyring       string
	endpoint      string
	secretKeyring string
	plaintext     bool
	machines      []string
)

func init() {
	flagset.StringVar(&backendName, "backend", "etcd", "backend provider")
	flagset.StringVar(&endpoint, "endpoint", "", "backend url")
	flagset.BoolVar(&plaintext, "plaintext", false, "skip encryption")
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
	case "list":
		listCmd(flagset)
	case "version":
		ver := buildVersion(Version, Commit, CommitDate, builtBy)
		versionCmd(flagset, ver)

	default:
		help()
	}
}

func help() {
	const usage = `usage: %s COMMAND [arg...]

commands:
   get   	retrieve the value of a key
   list  	retrieve all values under a key
   set   	set the value of a key
   version 	print the version of crypt

-plaintext  don't encrypt or decrypt the values before storage or retrieval
`

	_, _ = fmt.Fprintf(os.Stderr, usage, os.Args[0])
	os.Exit(1)
}

func buildVersion(version, commit, date, builtBy string) string {
	result := "crypt version " + version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	return result
}
