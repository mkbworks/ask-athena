package dns

const (
	ROOT_SERVER_ADDRESS = "8.8.8.8"
	DNS_PORT_NUMBER = 53
	MESSAGE_PROTOCOL = "udp"
	DOMAIN_LABEL_LIMIT = 63
	DOMAIN_NAME_LIMIT = 255
	DOMAIN_LABEL_SEPERATOR = "."
	UDP_MESSAGE_SIZE_LIMIT = 512
	MESSAGE_HEADER_LENGTH = 12
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

	RC_NO_ERROR ResponseCode = 0
	RC_FORMAT_ERROR ResponseCode = 1
	RC_SERVER_FAILURE ResponseCode = 2
	RC_NAME_ERROR ResponseCode = 3
	RC_NOT_IMPLEMENTED ResponseCode = 4
	RC_REFUSED ResponseCode = 5

	MSG_REQUEST MessageType = 201
	MSG_RESPONSE MessageType = 202
)