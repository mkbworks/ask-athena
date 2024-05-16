package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"path/filepath"
	"github.com/maheshkumaarbalaji/ask-athena/lib/dns"
)

func main() {
	t := flag.String("type", "A", "the record type to query for each domain name")
	flag.Parse()
	names := flag.Args()
	if len(names) == 0 {
		fmt.Println("Not enough arguments, must pass in at least one name")
		os.Exit(1)
	}
	CurrentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println("Error occurred while getting current working directory: " + err.Error())
		os.Exit(1)
	}

	CacheFilePath := filepath.Join(CurrentDirectory, "assets", "resolver-cache.conf")
	RootServersPath := filepath.Join(CurrentDirectory, "assets", "root-servers.conf")
	resolver, err := dns.GetResolver(RootServersPath, CacheFilePath)
	if err != nil {
		fmt.Printf("Error occurred while fetching DNS Resolver Instance: %s\n", err.Error())
		os.Exit(1)
	}

	if resolver.IsAllowed(*t) {
		for _, name := range names {
			values := resolver.Resolve(name, resolver.GetRecordType(*t))
			fmt.Printf("%s\t\t\t%s\n", name, strings.Join(values, "  "))
		}
	} else {
		fmt.Printf("Given record type is not supported by the DNS resolver.\n")
		os.Exit(1)
	}
}