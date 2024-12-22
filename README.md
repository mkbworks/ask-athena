# Ask Athena

<img src="https://mkbalaji.pages.dev/project-athena/project-athena.png" style="border-radius:50%" align="right" width="159px" alt="Project Athena logo">

A command-line based Recursive DNS resolver created using Golang. It is compliant with `RFC 1035` and supports the following record types.

- **A** record 
- **AAAA** record
- **CNAME** record 
- **TXT** record

The resolver also supports caching thereby facilitating quick resolution of domain names. The transfer of DNS messages, to and from the DNS server is done over User Datagram Protocol (UDP).

This is my solution to the challenge posted at [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-dns-resolver) to create my own DNS resolver.

## Example Usage

The `main.go` file in the root directory contains a sample code that can be used to invoke the DNS resolution process for a set of domain names given as command line arguments.

```go
recType := flag.String("type", "A", "the record type to query for each domain name")
traceLogs := flag.Bool("trace", false, "Enable/Disable Trace Logs")
flag.Parse()
names := flag.Args()
```

The above code snippet in main(), parses the command line arguments entered by the user. There are two types of arguments expected - a flag value to denote the type of DNS record being queried and whether or not trace logs are to be printed on screen, and a series of domain names for which the record type must be queried.

```go
config.SetupConfig()
```

The call to the SetupConfig() sets up the necessary configuration parameters required for creating an instance of the DNS resolver like Root DNS Servers file path and resolver cache file path.

```go
resolver, err := dns.NewResolver(config.RootServerFilePath, config.CacheFilePath, *traceLogs)
```

These are then used to create a new instance of the resolver by invoking the dns.NewResolver() method.

```go
resolver.Resolve(name, resolver.GetRecordType(*t))
```

Once the resolver is inititalized, the domain name resolution can be carried out by invoking the resolver.Resolve() function with each of the domain name and its record type. The GetRecordType() method gets the DNS record type object associated with the record type string fetched from the command line.

```go
resolver.Close()
```

Finally call the Close() method once all the domain names have been resolved. This persists the changes made to resolver cache in the local filesystem.

## Commands and Outputs

This section contains various examples of how `ask-athena` can leveraged to query for DNS records.

### Example 1

This example queries the `A` records (IPv4 address) for `www.scu.edu`

**Command entered** 

```bash
./ask-athena -type=A www.scu.edu
```

**Output printed on screen**

```bash
Querying DNS for A type record of www.scu.edu.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.scu.edu. 	 IN 	 A

ANSWER SECTION:
www.scu.edu. 	 30 	 IN 	 A 	 34.107.151.86

```

### Example 2

This example queries the `A` records (IPv4 address) for `www.mit.edu`

**Command entered** 

```bash
./ask-athena -type=A www.mit.edu
```

**Output printed on screen**

```bash
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

This example queries the `A` records (IPv4 address) for `www.facebook.com`

**Command entered** 

```bash
./ask-athena -type=A www.facebook.com
```

**Output printed on screen**

```bash
Querying DNS for A type record of www.facebook.com.

->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 0
Flags: QR RD RA, QUESTION: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
www.facebook.com. 	 IN 	 A

ANSWER SECTION:
www.facebook.com. 	 3600 	 IN 	 CNAME 	 star-mini.c10r.facebook.com.
star-mini.c10r.facebook.com. 	 60 	 IN 	 A 	 157.240.22.35

```

### Example 4

This example queries the `A` records (IPv4 address) for `google.com` with the `-trace` flag enabled to print the trace logs on screen. 

**Command entered** 

```bash
./ask-athena -type=A -trace=true google.com
```

**Output printed on screen**

```bash
Querying DNS for A type record of google.com.

2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 DNS Request being sent to server - 199.7.91.13.
2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Request Contents are:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: RD, QUESTION: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
google.com. 	 IN 	 A


2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Response received back:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: QR RD, QUESTION: 1, ANSWER: 0, AUTHORITY: 13, ADDITIONAL: 14

QUESTION SECTION:
google.com. 	 IN 	 A

AUTHORITY SECTION:
com. 	 172800 	 IN 	 NS 	 a.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 b.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 c.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 d.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 e.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 f.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 g.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 h.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 i.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 j.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 k.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 l.gtld-servers.net.
com. 	 172800 	 IN 	 NS 	 m.gtld-servers.net.

ADDITIONAL SECTION:
a.gtld-servers.net. 	 172800 	 IN 	 A 	 192.5.6.30
b.gtld-servers.net. 	 172800 	 IN 	 A 	 192.33.14.30
c.gtld-servers.net. 	 172800 	 IN 	 A 	 192.26.92.30
d.gtld-servers.net. 	 172800 	 IN 	 A 	 192.31.80.30
e.gtld-servers.net. 	 172800 	 IN 	 A 	 192.12.94.30
f.gtld-servers.net. 	 172800 	 IN 	 A 	 192.35.51.30
g.gtld-servers.net. 	 172800 	 IN 	 A 	 192.42.93.30
h.gtld-servers.net. 	 172800 	 IN 	 A 	 192.54.112.30
i.gtld-servers.net. 	 172800 	 IN 	 A 	 192.43.172.30
j.gtld-servers.net. 	 172800 	 IN 	 A 	 192.48.79.30
k.gtld-servers.net. 	 172800 	 IN 	 A 	 192.52.178.30
l.gtld-servers.net. 	 172800 	 IN 	 A 	 192.41.162.30
m.gtld-servers.net. 	 172800 	 IN 	 A 	 192.55.83.30
a.gtld-servers.net. 	 172800 	 IN 	 AAAA 	 2001:503:a83e::2:30

2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 DNS Request being sent to server - 192.5.6.30.
2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Request Contents are:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: RD, QUESTION: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
google.com. 	 IN 	 A


2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Response received back:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: QR RD, QUESTION: 1, ANSWER: 0, AUTHORITY: 4, ADDITIONAL: 8

QUESTION SECTION:
google.com. 	 IN 	 A

AUTHORITY SECTION:
google.com. 	 172800 	 IN 	 NS 	 ns2.google.com.
google.com. 	 172800 	 IN 	 NS 	 ns1.google.com.
google.com. 	 172800 	 IN 	 NS 	 ns3.google.com.
google.com. 	 172800 	 IN 	 NS 	 ns4.google.com.

ADDITIONAL SECTION:
ns2.google.com. 	 172800 	 IN 	 AAAA 	 2001:4860:4802:34::a
ns2.google.com. 	 172800 	 IN 	 A 	 216.239.34.10
ns1.google.com. 	 172800 	 IN 	 AAAA 	 2001:4860:4802:32::a
ns1.google.com. 	 172800 	 IN 	 A 	 216.239.32.10
ns3.google.com. 	 172800 	 IN 	 AAAA 	 2001:4860:4802:36::a
ns3.google.com. 	 172800 	 IN 	 A 	 216.239.36.10
ns4.google.com. 	 172800 	 IN 	 AAAA 	 2001:4860:4802:38::a
ns4.google.com. 	 172800 	 IN 	 A 	 216.239.38.10

2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 DNS Request being sent to server - 216.239.34.10.
2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Request Contents are:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: RD, QUESTION: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
google.com. 	 IN 	 A


2024/06/08 19:08:39 **********************************************
2024/06/08 19:08:39 Response received back:
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: QR AA RD, QUESTION: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
google.com. 	 IN 	 A

ANSWER SECTION:
google.com. 	 300 	 IN 	 A 	 142.251.32.46


2024/06/08 19:08:39 **********************************************
->> HEADER <<- Opcode: QUERY, Status: NOERROR, ID: 48339
Flags: QR RD RA, QUESTION: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

QUESTION SECTION:
google.com. 	 IN 	 A

ANSWER SECTION:
google.com. 	 300 	 IN 	 A 	 142.251.32.46

```