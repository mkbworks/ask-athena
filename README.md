# Ask Athena

<img src="https://mkbalaji.pages.dev/project-athena/project-athena.png" style="border-radius:50%" align="right" width="159px" alt="Project Athena logo">

A command-line based Recursive DNS resolver created using Golang. It is compliant with RFC 1035 and supports the following record types.

- **A** record type (IPv4 address mapped to the domain name)
- **AAAA** record type (IPv6 address mapped to the domain name)
- **CNAME** record type (Canonical name for the domain name)
- **TXT** record type (Text string configured for the domain name)

The resolver also supports caching thereby facilitating quick resolution of domain names. The transfer of DNS messages, to and from the DNS server is done over User Datagram Protocol (UDP).

## Example Usage

The **main.go** file in the root directory contains a sample code that can be used to invoke the DNS resolution process for a set of domain names given as command line arguments.

```text
t := flag.String("type", "A", "the record type to query for each domain name")
flag.Parse()
names := flag.Args()
```

The above code snippet in main(), parses the command line arguments entered by the user. There are two types of arguments expected - a flag value to denote the type of DNS record being queried, and a series of domain names for which the record type must be queried.

```text
config.SetupConfig(LogFileDirectory)
```

The call to the SetupConfig() sets up the necessary configuration parameters required for creating an instance of the DNS resolver like Root DNS Servers file path, resolver cache file path and log file path. These are then used to create a new instance of the resolver in the below snippet.

```text
resolver, err := dns.NewResolver(config.RootServerFilePath, config.CacheFilePath, config.LogFilePath)
```

Once the resolver is inititalized, the domain name resolution can be carried out by invoking the resolver.Resolve() function with each of the domain name and its record type.

```text
resolver.Resolve(name, resolver.GetRecordType(*t))
```

The GetRecordType() method gets the DNS record type object associated with the given record type string. Finally call the Close() method once all the entered domain names have been resolved. This persists the changes made to resolver cache from memory to the file present in local filesystem.

```text
resolver.Close()
```

## Example Output

Command entered: 

```text
./ask-athena -type=A www.scu.edu
```

Output printed on screen:

```text
Querying DNS for A type record of www.scu.edu.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.scu.edu. 	 IN 	 A

ANSWER SECTION:
www.scu.edu. 	 30 	 IN 	 A 	 34.107.151.86

```