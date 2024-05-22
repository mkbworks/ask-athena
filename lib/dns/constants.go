package dns

const (
	DNS_PORT_NUMBER = 53
	MESSAGE_PROTOCOL = "udp"
	DOMAIN_LABEL_SEPERATOR = "."
	UDP_MESSAGE_SIZE_LIMIT = 4096
	MESSAGE_HEADER_LENGTH = 12
	WHITESPACE = " "
	NEWLINE_SEPERATOR = "\n"
	ADDRESS_IPv4 = "IPv4"
	ADDRESS_IPv6 = "IPv6"
)

const (
	TYPE_A     RecordType = 1
	TYPE_NS    RecordType = 2
	TYPE_CNAME RecordType = 5
	TYPE_TXT   RecordType = 16
	TYPE_AAAA  RecordType = 28

	OPCODE_QUERY Flag = 0
	OPCODE_IQUERY Flag = 1
	OPCODE_STATUS Flag = 2

	CLASS_IN ClassType = 1
	CLASS_CH ClassType = 3

	// DNS Query completed successfully
	RC_NOERROR ResponseCode = 0
	// DNS Query Format Error
	RC_FORMERR ResponseCode = 1
	// Server failed to complete the DNS request
	RC_SERVFAIL ResponseCode = 2
	// Domain name does not exist
	RC_NXDOMAIN ResponseCode = 3
	// Function not implemented
	RC_NOTIMP ResponseCode = 4
	// The server refused to answer for the query
	RC_REFUSED ResponseCode = 5
	// Name that should not exist, does exist
	RC_YXDOMAIN ResponseCode = 6
	// RRset that should not exist, does exist
	RC_XRRSET ResponseCode = 7
	// Server not authoritative for the zone
	RC_NOTAUTH ResponseCode = 8
	// Name not in zone
	RC_NOTZONE ResponseCode = 9

	MSG_REQUEST MessageType = 0
	MSG_RESPONSE MessageType = 1
	MSG_RESOLVER_RESPONSE MessageType = 2

	QR_BIT = uint16(1 << 15)
	AA_BIT = uint16(1 << 10)
	TR_BIT = uint16(1 << 9)
	RD_BIT = uint16(1 << 8)
	RA_BIT = uint16(1 << 7)
	AUTH_BIT = uint16(1 << 5)
	CHK_BIT = uint16(1 << 4)
	RCODE_BITS = uint16(15)
	OPCODE_BITS = uint16(15 << 11)
	PTR_DETECT_VALUE = uint16(3 << 14)
	PTR_OFFSET_FETCH = uint16(65535 >> 2)
)

var AllowedRRTypes RecordTypes = RecordTypes{
	"A":     TYPE_A,
	"NS":    TYPE_NS,
	"CNAME": TYPE_CNAME,
	"TXT":   TYPE_TXT,
	"AAAA":  TYPE_AAAA,
}

var AllowedClassTypes ClassTypes = ClassTypes{
	"IN": CLASS_IN,
	"CH": CLASS_CH,
}