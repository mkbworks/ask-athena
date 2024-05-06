package dns

import "strconv"

//Represents the multi-bit flags to be configured in DNS message header.
type Flag uint16
//Represents the response code returned by the DNS server.
type ResponseCode uint16

//Represents the header of a DNS message.
type Header struct {
	//Unique Identifier assigned to the DNS Request.
	Identifier uint16
	//Indicates if the DNS Message object is a request or response. The QR flag is set to 1, if value is true and is set to 0, if value is false.
	IsResponse bool
	Opcode Flag
	Authoritative bool
	Truncation bool
	RecursionDesired bool
	RecursionAvailable bool
	Zero Flag
	Rcode ResponseCode
	QdCount uint16
	AnCount uint16
	NsCount uint16
	ArCount uint16
}

//Set the default values for both request and response headers in a DNS Message.
func (hdr *Header) SetDefaults(mt MessageType) {
	hdr.Identifier = Id()
	hdr.Zero = 0
	if mt == MSG_REQUEST {
		hdr.IsResponse = false
	} else {
		hdr.IsResponse = true
	}
	
	hdr.Opcode = OPCODE_QUERY
	hdr.Authoritative = false
	hdr.Truncation = false
	hdr.RecursionDesired = false
	hdr.RecursionAvailable = false
	hdr.Rcode = RC_NO_ERROR
	hdr.QdCount = 0
	hdr.AnCount = 0
	hdr.ArCount = 0
	hdr.NsCount = 0
}

//Sets the number of questions in the DNS Message.
func (hdr *Header) SetQuestionCount(count uint16) {
	hdr.QdCount = count
}

//Sets the number of answer RRs in the DNS Message.
func (hdr *Header) SetAnswerCount(count uint16) {
	hdr.AnCount = count
}

//Sets the number of authoritative RRs present in the DNS Message.
func (hdr *Header) SetNameServerCount(count uint16) {
	hdr.NsCount = count
}

//Sets the number of Additional RRs present in the DNS Message.
func (hdr *Header) SetAdditionalRecordCount(count uint16) {
	hdr.ArCount = count
}

//Pack the flag values in Message Header as a binary string.
func (hdr *Header) PackFlag() []byte {
	flag_binary := ""
	if hdr.IsResponse {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	flag_binary += GetBinary(uint16(hdr.Opcode), 4)

	if hdr.Authoritative {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	if hdr.Truncation {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	if hdr.RecursionDesired {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	if hdr.RecursionAvailable {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	flag_binary += GetBinary(uint16(hdr.Zero), 3)
	flag_binary += GetBinary(uint16(hdr.Rcode), 4)
	return PackBinary16(flag_binary)
}

//Pack the header instance as a sequence of octets.
func (hdr *Header) Pack() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, PackUInt16(hdr.Identifier)...)
	buffer = append(buffer, hdr.PackFlag()...)
	buffer = append(buffer, PackUInt16(hdr.QdCount)...)
	buffer = append(buffer, PackUInt16(hdr.AnCount)...)
	buffer = append(buffer, PackUInt16(hdr.NsCount)...)
	buffer = append(buffer, PackUInt16(hdr.ArCount)...)
	return buffer
}

//Unpacks a flag byte stream into the Header instance.
func (hdr *Header) UnpackFlag(buffer []byte) {
	binary_string := UnpackBinary16(buffer)
	QrBit := binary_string[:1]
	if QrBit == "1" {
		hdr.IsResponse = true
	} else {
		hdr.IsResponse = false
	}

	OpcodeBits := binary_string[1:5]
	OpcodeValue, err := strconv.ParseUint(OpcodeBits, 2, 16)
	if err != nil {
		panic(err)
	}

	hdr.Opcode = Flag(OpcodeValue)
	AABit := binary_string[5:6]
	if AABit == "1" {
		hdr.Authoritative = true
	} else {
		hdr.Authoritative = false
	}

	TCBit := binary_string[6:7]
	if TCBit == "1" {
		hdr.Truncation = true
	} else {
		hdr.Truncation = false
	}

	RDBit := binary_string[7:8]
	if RDBit == "1" {
		hdr.RecursionDesired = true
	} else {
		hdr.RecursionDesired = false
	}

	RABit := binary_string[8:9]
	if RABit == "1" {
		hdr.RecursionAvailable = true
	} else {
		hdr.RecursionAvailable = false
	}

	ZeroBits := binary_string[9:12]
	ZeroValue, err := strconv.ParseUint(ZeroBits, 2, 16)
	if err != nil {
		panic(err)
	}

	hdr.Zero = Flag(ZeroValue)
	RCodeBits := binary_string[12:16]
	RCodeValue, err := strconv.ParseUint(RCodeBits, 2, 16)
	if err != nil {
		panic(err)
	}

	hdr.Rcode = ResponseCode(RCodeValue)
}

//Unpacks a stream of bytes to a header instance of DNS Message.
func (hdr *Header) Unpack(buffer []byte) {
	hdr.Identifier = UnpackUInt16(buffer[:2])
	hdr.UnpackFlag(buffer[2:4])
	hdr.QdCount = UnpackUInt16(buffer[4:6])
	hdr.AnCount = UnpackUInt16(buffer[6:8])
	hdr.NsCount = UnpackUInt16(buffer[8:10])
	hdr.ArCount = UnpackUInt16(buffer[10:])
}