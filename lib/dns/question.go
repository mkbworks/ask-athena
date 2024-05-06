package dns

//Represents a DNS Question record.
type Question struct {
	Name DomainName
	Type RecordType
	Class ClassType
}

//Sets the Question instance with the given values.
func (que *Question) Set(name string, recType RecordType) {
	que.Class = CLASS_IN
	que.Type = recType
	que.Name = DomainName{}
	que.Name.Encode(name)
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
func (que *Question) Unpack(buffer []byte) {

}