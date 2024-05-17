package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"github.com/maheshkumaarbalaji/ask-athena/lib/dns"
	"github.com/maheshkumaarbalaji/ask-athena/lib/config"
)

func main() {
	t := flag.String("type", "A", "the record type to query for each domain name")
	flag.Parse()
	names := flag.Args()
	if len(names) == 0 {
		fmt.Println("Not enough arguments, must pass in at least one name")
		os.Exit(1)
	}

	err := config.SetupConfig()
	if err != nil {
		fmt.Println("Error occurred while setting up DNS resolver configuration:", err.Error())
		os.Exit(1)
	}

	resolver, err := dns.GetResolver(config.RootServerFilePath, config.CacheFilePath, config.LogFilePath)
	if err != nil {
		fmt.Printf("Error occurred while fetching DNS Resolver Instance: %s\n", err.Error())
		os.Exit(1)
	}

	if resolver.IsAllowed(*t) {
		for _, name := range names {
			values := resolver.Resolve(name, resolver.GetRecordType(*t))
			fmt.Printf("%s\t\t%s\n", name, strings.Join(values, ", "))
		}
	} else {
		fmt.Printf("Given record type is not supported by the DNS resolver.\n")
	}
	resolver.Close()
}