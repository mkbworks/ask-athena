package dns

import (
	"bytes"
	"errors"
	"strings"
)

//Represents the type of DNS Message - Request or Response.
type MessageType uint8

//Represents a domain name as defined in RFC 1035.
type DomainName struct {
	Data []byte
	Value string
	Length uint8
}

//Parse the given domain name string and store it as sequence of octets.
func (name *DomainName) Encode(dName string) {
	encodedBytes := make([]byte, 0)
	dName = strings.TrimSuffix(dName, DOMAIN_LABEL_SEPERATOR)
	labelCount := 0
	for _, dLabel := range strings.Split(dName, DOMAIN_LABEL_SEPERATOR) {
		labelBytes := []byte(dLabel)
		encodedBytes = append(encodedBytes, byte(len(labelBytes)))
		encodedBytes = append(encodedBytes, labelBytes...)
		labelCount++
	}

	encodedBytes = append(encodedBytes, byte(0))
	name.Data = encodedBytes
	name.Length = uint8(labelCount)
	name.Value = dName
}

//Decode the given byte stream and extract the domain name.
func (name *DomainName) Decode(buffer []byte) {
	name.Data = buffer
	completeDomainName := ""
	LastIndexRead := 0
	labelByteCount := buffer[LastIndexRead]
	name.Length = 0
	for iterate := true; iterate; {
		if int(labelByteCount) != 0 {
			labelBytes := buffer[LastIndexRead + 1: LastIndexRead + int(labelByteCount) + 1]
			completeDomainName += string(labelBytes) + DOMAIN_LABEL_SEPERATOR
			LastIndexRead = LastIndexRead + int(labelByteCount)
			name.Length++
		} else {
			iterate = false
		}
	}

	name.Value = strings.TrimSuffix(completeDomainName, DOMAIN_LABEL_SEPERATOR)
}

//Represents a DNS Message (both Request and Response).
type Message struct {
	Header Header
	Questions []Question
	Answers []Resource
	Authoritative []Resource
	Additional []Resource
}

//Set Request values for the Message instance.
func (msg *Message) SetRequest(name string, recType RecordType) {
	msg.Header.SetDefaults(MSG_REQUEST)
	question := Question{}
	question.Set(name, recType)
	msg.Questions = append(msg.Questions, question)
	msg.Header.SetQuestionCount(uint16(len(msg.Questions)))
}

//Pack the message as a sequence of octets.
func (msg *Message) Pack() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, msg.Header.Pack()...)
	if msg.Header.QdCount > 0 {
		for _, que := range msg.Questions {
			buffer = append(buffer, que.Pack()...)
		}
	}
	return buffer
}

//Unpack the sequence of bytes to a Message instance.
func (msg *Message) Unpack(response []byte) {
	buffer := make([]byte, len(response))
	copy(buffer, response)
	msg.Header.Unpack(buffer[:MESSAGE_HEADER_LENGTH])
	buffer = buffer[MESSAGE_HEADER_LENGTH:]
	if msg.Header.QdCount > 0 {
		for index := 1; index <= int(msg.Header.QdCount); index++ {
			indexOne := bytes.IndexByte(buffer, byte(0))
			if indexOne != -1 {
				temp := buffer[:indexOne + 1]
				question := Question{}
				question.Unpack(temp)
				msg.Questions = append(msg.Questions, question)
			} else {
				panic(errors.New("domain name delimiter byte missing in given byte stream"))
			}
		}
	}
}