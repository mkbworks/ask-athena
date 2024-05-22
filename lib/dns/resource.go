package dns

import (
	"fmt"
	"strings"
)

//Feature(s) to be implemented for a DNS Resource Body.
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

//Initialize the instance of Resource.
func (resource *Resource) Initialize(domainName string, recType string, classType string, ttl uint32, data string) {
	resource.Name = DomainName{}
	resource.Name.Initialize(domainName)
	resource.Type = AllowedRRTypes.GetRecordType(recType)
	resource.Class = AllowedClassTypes.GetClassType(classType)
	resource.TTL = ttl
	if resource.Type == TYPE_A {
		resource.RdLength = uint16(len(convertToBytes(data, ADDRESS_IPv4)))
		ar := AResource{}
		ar.IPv4Address = strings.TrimSpace(data)
		resource.Rdata = &ar
	} else if resource.Type == TYPE_AAAA {
		resource.RdLength = uint16(len(convertToBytes(data, ADDRESS_IPv6)))
		aaar := AAAAResource{}
		aaar.IPv6Address = strings.TrimSpace(data)
		resource.Rdata = &aaar
	} else if resource.Type == TYPE_CNAME {
		cname := CNAMEResource{}
		cname.name = DomainName{}
		cname.name.Initialize(data)
		resource.RdLength = uint16(cname.name.GetLength())
		resource.Rdata = &cname
	} else if resource.Type == TYPE_NS {
		ns := NSResource{}
		ns.NameServer = DomainName{}
		ns.NameServer.Initialize(data)
		resource.RdLength = uint16(ns.NameServer.GetLength())
		resource.Rdata = &ns
	} else if resource.Type == TYPE_TXT {
		txt := TXTResource{}
		txt.TextValue = strings.TrimSpace(data)
		resource.RdLength = uint16(len(data))
		resource.Rdata = &txt
	}
}

//Packs the resource data into a stream of bytes
func (resource *Resource) PackBody(compressionMap CompressionMap, offset int) []byte {
	buffer := make([]byte, 0)
	if resource.Type == TYPE_A {
		resourceData := resource.GetData()
		buffer = append(buffer, convertToBytes(resourceData,  ADDRESS_IPv4)...)
	} else if resource.Type == TYPE_AAAA {
		resourceData := resource.GetData()
		buffer = append(buffer, convertToBytes(resourceData,  ADDRESS_IPv6)...)
	} else if resource.Type == TYPE_CNAME || resource.Type == TYPE_NS {
		buffer = append(buffer, resource.Name.Pack(compressionMap, offset)...)
	} else if resource.Type == TYPE_TXT {
		buffer = append(buffer, []byte(resource.GetData())...)
	}

	return buffer
}

//Packs the resource instance to a stream of bytes
func (resource *Resource) Pack(compressionMap CompressionMap, offset int) []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, resource.Name.Pack(compressionMap, offset)...)
	buffer = append(buffer, PackUInt16(uint16(resource.Type))...)
	buffer = append(buffer, PackUInt16(uint16(resource.Class))...)
	buffer = append(buffer, PackUInt32(uint32(resource.TTL))...)
	buffer = append(buffer, PackUInt16(resource.RdLength)...)
	offset = offset + len(buffer)
	buffer = append(buffer, resource.PackBody(compressionMap, offset)...)
	return buffer
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
	return fmt.Sprintf("%s \t %d \t %s \t %s \t %s\n", resource.Name.String(), int(resource.TTL) , resource.Class.String(), resource.Type.String(), value_string)
}

//Returns a string representation of the Resource instance for caching
func (resource *Resource) CacheString() string {
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
	values := make([]string, 0)
	values = append(values, resource.Name.String())
	values = append(values, fmt.Sprintf("%d", int(resource.TTL)))
	values = append(values, resource.Class.String())
	values = append(values, resource.Type.String())
	values = append(values, value_string)

	return strings.Join(values, WHITESPACE)
}


//Gets the value in the resource body.
func (resource *Resource) GetData() string {
	value_string := ""
	if obj, ok := resource.Rdata.(*AResource); ok {
		value_string = obj.IPv4Address
	} else if obj, ok := resource.Rdata.(*AAAAResource); ok {
		value_string = obj.IPv6Address
	} else if obj, ok := resource.Rdata.(*CNAMEResource); ok {
		value_string = obj.name.Value
	} else if obj, ok := resource.Rdata.(*NSResource); ok {
		value_string = obj.NameServer.Value
	} else if obj, ok := resource.Rdata.(*TXTResource); ok {
		value_string = obj.TextValue
	} else {
		value_string = ""
	}
	return value_string
}

//Represents an A-type Resource Record value.
type AResource struct {
	IPv4Address string
}

//Unpacks a stream of bytes into a A-type resource record value.
func (ar *AResource) UnpackBody(buffer []byte, offset int, dataLength int) int {
	ipBytes := buffer[offset: offset + dataLength]
	ar.IPv4Address = getIPAddress(ipBytes)
	return offset + dataLength
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
	aaaar.IPv6Address = getIPAddress(ipBytes)
	return offset + dataLength
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