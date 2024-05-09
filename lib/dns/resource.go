package dns

import (
	"fmt"
	"net"
)

//Contains all the functions to implement for parsing a DNS message body.
type ResourceBody interface {
	//Unpacks a stream of bytes into a resource record object.
	UnpackBody(buffer []byte, offset int, dataLength int) int
}

//Represents a Resource Record in DNS.
type Resource struct {
	//Contains the domain name value being returned.
	Name DomainName
	//Contains the type of record represented by the RR.
	Type RecordType
	//Contains the class type represented by the RR.
	Class ClassType
	//Time-To-Live for the RR.
	TTL uint32
	//Total length (in bytes) of the Resource body.
	RdLength uint16
	//Resource Record body.
	Rdata ResourceBody
}

//Unpacks a stream of bytes to a resource instance.
func (resource *Resource) Unpack(buffer []byte, offset int) int {
	offset = resource.Name.Unpack(buffer, offset)
	resource.Type = RecordType(UnpackUInt16(buffer[offset: offset + 2]))
	resource.Class = ClassType(UnpackUInt16(buffer[offset + 2: offset + 4]))
	resource.TTL = UnpackUInt32(buffer[offset + 4: offset + 8])
	resource.RdLength = UnpackUInt16(buffer[offset + 8: offset + 10])
	if resource.Type == TYPE_A {
		ar := AResource{}
		offset = ar.UnpackBody(buffer, offset + 10, int(resource.RdLength))
		resource.Rdata = &ar
	} else if resource.Type == TYPE_AAAA {
		aaar := AAAAResource{}
		offset = aaar.UnpackBody(buffer, offset + 10, int(resource.RdLength))
		resource.Rdata = &aaar
	} else if resource.Type == TYPE_CNAME {
		cname := CNAMEResource{}
		offset = cname.UnpackBody(buffer, offset + 10, int(resource.RdLength))
		resource.Rdata = &cname
	} else if resource.Type == TYPE_NS {
		ns := NSResource{}
		offset = ns.UnpackBody(buffer, offset + 10, int(resource.RdLength))
		resource.Rdata = &ns
	} else if resource.Type == TYPE_TXT {
		txt := TXTResource{}
		offset = txt.UnpackBody(buffer, offset + 10, int(resource.RdLength))
		resource.Rdata = &txt
	}

	return offset
}

//Returns a string representation of the Resource instance.
func (resource *Resource) String() string {
	value_string := ""
	if obj, ok := resource.Rdata.(*AResource); ok {
		value_string = obj.String()
	} else if obj, ok := resource.Rdata.(*AAAAResource); ok {
		value_string = obj.String()
	} else if obj, ok := resource.Rdata.(*CNAMEResource); ok {
		value_string = obj.String()
	} else if obj, ok := resource.Rdata.(*NSResource); ok {
		value_string = obj.String()
	} else if obj, ok := resource.Rdata.(*TXTResource); ok {
		value_string = obj.String()
	} else {
		value_string = ""
	}
	return fmt.Sprintf("%s \t %s \t %s \t %d \t %s\n", resource.Name.String(), resource.Class.String(), resource.Type.String(), int(resource.TTL), value_string)
}

//Represents an A-type Resource Record value.
type AResource struct {
	IPv4Address string
}

//Unpacks a stream of bytes into a A-type resource record value.
func (ar *AResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	ipBytes := buffer[offset: offset + dataLength]
	ipAddress := net.IP(ipBytes)
	ar.IPv4Address = ipAddress.String()
	return offset
}

//Returns the string representation of A-type record data
func (ar *AResource) String() string {
	return ar.IPv4Address
}

//Represents a AAAA-type Resource Record body.
type AAAAResource struct {
	IPv6Address string
}

//Unpacks a stream of bytes into a AAAA-type resource record value.
func (aaaar *AAAAResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	ipBytes := buffer[offset: offset + dataLength]
	ipAddress := net.IP(ipBytes)
	aaaar.IPv6Address = ipAddress.String()
	return offset
}

//Returns the string representation of AAAA-type data.
func (aaar *AAAAResource) String() string {
	return aaar.IPv6Address
}

//Represents a CNAME-type Resouce Record body.
type CNAMEResource struct {
	name DomainName
}

//Unpacks a stream of bytes into a CNAME-type resource record value.
func (cname *CNAMEResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	cname.name = DomainName{}
	offset = cname.name.Unpack(buffer, offset)
	return offset
}

//Returns the string representation of CNAME-type record value.
func (cname *CNAMEResource) String() string {
	return cname.name.String()
}

//Represents a NS-type Resource Record body.
type NSResource struct {
	NameServer DomainName
}

//Unpacks a stream of bytes into a NS-type resource record value.
func (ns *NSResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	offset = ns.NameServer.Unpack(buffer, offset)
	return offset
}

//Returns the string representation of NS-type record value.
func (ns *NSResource) String() string {
	return ns.NameServer.String()
}

//Represents a TXT-type Resource Record body
type TXTResource struct {
	TextValue string
}

//Unpacks a stream of bytes into a TXT-type resource record value.
func (txt *TXTResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	txtByteSlice := buffer[offset: offset + int(dataLength)]
	offset = offset + int(dataLength)
	txt.TextValue = string(txtByteSlice)
	return offset
}

//Returns the TXT value.
func (txt *TXTResource) String() string {
	return txt.TextValue
}