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

Once the resolver is inititalized, the domain name resolution can be carried out by invoking the resolver.Resolve() function with each of the domain name and its record type. The GetRecordType() method gets the DNS record type object associated with the record type string fetched from the command line.

```text
resolver.Resolve(name, resolver.GetRecordType(*t))
```

Finally call the Close() method once all the domain names have been resolved. This persists the changes made to resolver cache in the local filesystem.

```text
resolver.Close()
```

## Commands and Outputs

### Example 1

**Command entered** 

```text
./ask-athena -type=A www.scu.edu
```

**Output printed on screen**

```text
Querying DNS for A type record of www.scu.edu.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.scu.edu. 	 IN 	 A

ANSWER SECTION:
www.scu.edu. 	 30 	 IN 	 A 	 34.107.151.86

```

### Example 2

**Command entered** 

```text
./ask-athena -type=A www.mit.edu
```

**Output printed on screen**

```text
Querying DNS for A type record of www.mit.edu.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.mit.edu. 	 IN 	 A

ANSWER SECTION:
www.mit.edu. 	 1800 	 IN 	 CNAME 	 www.mit.edu.edgekey.net.
www.mit.edu.edgekey.net. 	 60 	 IN 	 CNAME 	 e9566.dscb.akamaiedge.net.
e9566.dscb.akamaiedge.net. 	 20 	 IN 	 A 	 23.203.236.99

```

### Example 3

**Command entered** 

```text
./ask-athena -type=A www.facebook.com
```

**Output printed on screen**

```text
Querying DNS for A type record of www.facebook.com.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.facebook.com. 	 IN 	 A

ANSWER SECTION:
www.facebook.com. 	 3600 	 IN 	 CNAME 	 star-mini.c10r.facebook.com.
star-mini.c10r.facebook.com. 	 60 	 IN 	 A 	 157.240.22.35

```
