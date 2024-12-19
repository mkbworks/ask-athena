package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/mkbworks/ask-athena/lib/dns"
	"github.com/mkbworks/ask-athena/lib/config"
)

func main() {
	recType := flag.String("type", "A", "the record type to query for each domain name")
	traceLogs := flag.Bool("trace", false, "Enable/Disable Trace Logs")
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

	resolver, err := dns.NewResolver(config.RootServerFilePath, config.CacheFilePath, *traceLogs)
	if err != nil {
		fmt.Printf("Error occurred while fetching DNS Resolver Instance: %s\n", err.Error())
		os.Exit(1)
	}

	if resolver.IsAllowed(*recType) {
		for _, name := range names {
			fmt.Printf("Querying DNS for %s type record of %s.\n\n", *recType, name)
			resolver.Resolve(name, resolver.GetRecordType(*recType))
		}
	} else {
		fmt.Printf("Given record type is not supported by the DNS resolver.\n")
	}
	resolver.Close()
}