package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"github.com/maheshkumaarbalaji/project-athena/lib/dns"
)

func main() {
	t := flag.String("type", "A", "the record type to query for each domain name")
	flag.Parse()
	names := flag.Args()
	fmt.Printf("Domain name(s) to be resolved: %s\n", strings.Join(names, " , "))
	if len(names) == 0 {
		fmt.Println("Not enough arguments, must pass in at least one name")
		os.Exit(1)
	}

	resolver, err := dns.GetResolver()
	if err != nil {
		fmt.Printf("Error occurred while fetching DNS Resolver Instance: %s\n", err.Error())
		os.Exit(1)
	}

	if resolver.IsAllowed(*t) {
		for _, name := range names {
			resolver.Resolve(name, resolver.GetRecordType(*t))
		}
	} else {
		fmt.Printf("Given record type is not supported by the DNS resolver.\n")
		os.Exit(1)
	}

	resolver.Close()
}