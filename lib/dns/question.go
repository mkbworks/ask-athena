package dns

import (
	"fmt"
)

//Represents a DNS Question record.
type Question struct {
	//Contains the domain name being queried to the DNS Server.
	Name DomainName
	//Contains the type of record being queried to the DNS Server.
	Type RecordType
	//Contains the class type of the record being queried to the DNS Server.
	Class ClassType
}

//Sets the Question instance with the given values.
func (que *Question) Set(name string, recType RecordType) {
	que.Class = CLASS_IN
	que.Type = recType
	que.Name = DomainName{}
	que.Name.Initialize()
	que.Name.Pack(name)
}

//Pack the Question instance as a sequence of octets.
func (que *Question) Pack() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, que.Name.Data...)
	buffer = append(buffer, PackUInt16(uint16(que.Type))...)
	buffer = append(buffer, PackUInt16(uint16(que.Class))...)
	return buffer
}

//Unpacks a stream of bytes to a Question instance.
func (que *Question) Unpack(buffer []byte, offset int) int {
	offset = que.Name.Unpack(buffer, offset)
	que.Type = RecordType(UnpackUInt16(buffer[offset: offset + 2]))
	que.Class = ClassType(UnpackUInt16(buffer[offset + 2: offset + 4]))
	return offset + 4
}

//Returns the string representation of DNS Question instance.
func (que *Question) String() string {
	return fmt.Sprintf("%s \t %s \t %s\n", que.Name.String(), que.Class.String(), que.Type.String())
}