# Ask Athena

<img src="https://mkbalaji.pages.dev/project-athena/project-athena.png" style="border-radius:50%" align="right" width="159px" alt="Project Athena logo">

A command-line based Recursive DNS resolver created using Golang. It is compliant with RFC 1035 and supports the following record types.

- **A** record type (IPv4 address mapped to the domain name)
- **AAAA** record type (IPv6 address mapped to the domain name)
- **CNAME** record type (Canonical name for the domain name)
- **TXT** record type (Text string configured for the domain name)

The resolver also supports caching thereby facilitating quick resolution of domain names. The transfer of DNS messages, to and from the DNS server is done over User Datagram Protocol (UDP).
