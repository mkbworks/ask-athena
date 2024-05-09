package dns

import (
	"fmt"
	"strconv"
	"strings"
)

//Represents the multi-bit flags to be configured in DNS message header.
type Flag uint16

//Returns the string representation of the flag value.
func (flg Flag) String() string {
	switch flg {
	case OPCODE_QUERY:
		return "QUERY"
	case OPCODE_IQUERY:
		return "IQUERY"
	case OPCODE_STATUS:
		return "STATUS"
	default:
		return "" 
	}
}

//Represents the response code returned by the DNS server.
type ResponseCode uint16

//Returns the string representation of the response code.
func (rc ResponseCode) String() string {
	switch rc {
	case RC_NOERROR:
		return "NOERROR"
	case RC_FORMERR:
		return "FORMERR"
	case RC_SERVFAIL:
		return "SERVFAIL"
	case RC_NXDOMAIN:
		return "NXDOMAIN"
	case RC_NOTIMP:
		return "NOTIMP"
	case RC_REFUSED:
		return "REFUSED"
	case RC_YXDOMAIN:
		return "YXDOMAIN"
	case RC_XRRSET:
		return "XRRSET"
	case RC_NOTAUTH:
		return "NOTAUTH"
	case RC_NOTZONE:
		return "NOTZONE"
	default:
		return ""
	}
}

//Represents the header of a DNS message.
type Header struct {
	//Unique Identifier assigned to the DNS Request.
	Identifier uint16
	//Indicates if the DNS Message object is a request or response.
	IsResponse bool
	//Represents the type of query - Standard query or IQuery.
	Opcode Flag
	//Indicates if the RRs were returned by an authoritative server
	Authoritative bool
	//Indicates if the RRs in the DNS response are truncated.
	Truncation bool
	//Indicates if the resolver requires the target DNS server to fetch results through recursive name queries
	RecursionDesired bool
	//Indicates if the DNS server is capable of performing recursive name queries.
	RecursionAvailable bool
	//1-bit code reserved for future use.
	Zero Flag
	//1-bit code to indicate if data present after DNS header has been authenticated by the server.
	Authenticated bool
	//1-bit code to indicate if checking is disabled by the DNS Resolver.
	CheckingDisabled bool
	//Response code
	Rcode ResponseCode
	//Number of Question records in the DNS message
	QdCount uint16
	//Number of Answer records in the DNS message
	AnCount uint16
	//Number of Authoritative RRs in the DNS message
	NsCount uint16
	//Number of additional records present in the DNS message
	ArCount uint16
}

//Initialises an instance of Header with the default values based on MessageType.
func (hdr *Header) Initialize(mt MessageType) {
	if mt == MSG_REQUEST {
		hdr.Identifier = Id()
		hdr.IsResponse = false
	} else {
		hdr.Identifier = 0
		hdr.IsResponse = true
	}

	hdr.Opcode = OPCODE_QUERY
	hdr.Rcode = RC_NOERROR
	hdr.Zero = 0
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

	flag_binary += "0" //Default value for zero flag

	if hdr.Authenticated {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

	if hdr.CheckingDisabled {
		flag_binary += "1"
	} else {
		flag_binary += "0"
	}

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
	binary_string := UnpackBinary(buffer)
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

	hdr.Zero = Flag(0)
	AuthBit := binary_string[10:11]
	if AuthBit == "1" {
		hdr.Authenticated = true
	} else {
		hdr.Authenticated = false
	}

	ChkBit := binary_string[11:12]
	if ChkBit == "1" {
		hdr.CheckingDisabled = true
	} else {
		hdr.CheckingDisabled = false
	}
	
	RCodeBits := binary_string[12:16]
	RCodeValue, err := strconv.ParseUint(RCodeBits, 2, 16)
	if err != nil {
		panic(err)
	}

	hdr.Rcode = ResponseCode(RCodeValue)
}

//Unpacks a stream of bytes to a header instance of DNS Message.
func (hdr *Header) Unpack(buffer []byte, offset int) int {
	hdr.Identifier = UnpackUInt16(buffer[offset: offset + 2])
	hdr.UnpackFlag(buffer[offset + 2:offset + 4])
	hdr.SetQuestionCount(UnpackUInt16(buffer[offset + 4: offset + 6]))
	hdr.SetAnswerCount(UnpackUInt16(buffer[offset + 6: offset + 8])) 
	hdr.SetNameServerCount(UnpackUInt16(buffer[offset + 8: offset + 10]))
	hdr.SetAdditionalRecordCount(UnpackUInt16(buffer[offset + 10: offset + 12]))
	return offset + 12
}

//Returns the string representation of DNS Message Header.
func (hdr *Header) String() string {
	return_value := fmt.Sprintf("->> HEADER <<- Opcode: %s, Status: %s, ID: %d\n", hdr.Opcode.String(), hdr.Rcode.String(), int(hdr.Identifier))
	flags := make([]string, 0)
	if hdr.IsResponse {
		flags = append(flags, "QR")
	}

	if hdr.Authoritative {
		flags = append(flags, "AA")
	}

	if hdr.Truncation {
		flags = append(flags, "TC")
	}

	if hdr.RecursionDesired {
		flags = append(flags, "RD")
	}

	if hdr.RecursionAvailable {
		flags = append(flags, "RA")
	}

	return_value += fmt.Sprintf("Flags: %s, QUESTION: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n", strings.Join(flags, " "), int(hdr.QdCount), int(hdr.AnCount), int(hdr.NsCount), int(hdr.ArCount))
	return return_value
}