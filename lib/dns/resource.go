package dns

//Represents a Resource Record in DNS.
type Resource struct {
	Name DomainName
	Type RecordType
	Class ClassType
	TTL uint32
	RdLength uint16
	Rdata []byte
}