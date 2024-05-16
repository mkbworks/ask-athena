package dns

const (
	DNS_PORT_NUMBER = 53
	MESSAGE_PROTOCOL = "udp"
	DOMAIN_LABEL_LIMIT = 63
	DOMAIN_NAME_LIMIT = 255
	DOMAIN_LABEL_SEPERATOR = "."
	UDP_MESSAGE_SIZE_LIMIT = 512
	MESSAGE_HEADER_LENGTH = 12
	DOMAIN_NAME_PTR_PREFIX = "11"
	WHITESPACE = " "
	NEWLINE_SEPERATOR = "\n"
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

	RC_NOERROR ResponseCode = 0
	RC_FORMERR ResponseCode = 1
	RC_SERVFAIL ResponseCode = 2
	RC_NXDOMAIN ResponseCode = 3
	RC_NOTIMP ResponseCode = 4
	RC_REFUSED ResponseCode = 5
	RC_YXDOMAIN ResponseCode = 6
	RC_XRRSET ResponseCode = 7
	RC_NOTAUTH ResponseCode = 8
	RC_NOTZONE ResponseCode = 9

	MSG_REQUEST MessageType = 201
	MSG_RESPONSE MessageType = 202

	ROOT_SERVER ServerType = 1
	TLD_SERVER ServerType = 2
	AUTHORITATIVE_SERVER ServerType = 3
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