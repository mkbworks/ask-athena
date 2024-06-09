package dns

import (
	"strings"
)

//Represents the type of DNS Message - Request or Response.
type MessageType uint8

//Holds the domain name and their compression offsets used while packing them into a stream of octets.
type CompressionMap map[string]int

//Gets the offset assosciated with a domain name from the compression map.
func (cmpMap CompressionMap) Get(name string) (int, bool) {
	name = Canonicalize(name)
	offset, ok := cmpMap[name]
	return offset, ok
}

//Sets the offset for a domain name in the compression map. If domain name does not exist already, a new entry is added to the compression map.
func (cmpMap CompressionMap) Set(name string, offset int) {
	name = Canonicalize(name)
	cmpMap[name] = offset
}

//Represents a DNS Message (both Request and Response).
type Message struct {
	//Represents all the data present in header section of the DNS Message
	Header Header
	//Array of Question records in the DNS Message
	Questions []Question
	//Array of answer RRs present in the DNS Message.
	Answers []Resource
	//Array of authoritative RRs present in the DNS Message.
	Authoritative []Resource
	//Array of additional RRs present in the DNS Message.
	Additional []Resource
	//Compression map instance for the Message.
	compressionMap CompressionMap
}

//Initialises all the properties in the Message instance.
func (msg *Message) Initialize(mt MessageType, MsgId uint16) {
	msg.Header = Header{}
	msg.Header.Initialize(mt, MsgId)
	msg. Questions = make([]Question, 0)
	msg.Answers = make([]Resource, 0)
	msg.Authoritative = make([]Resource, 0)
	msg.Additional = make([]Resource, 0)
	msg.compressionMap = make(CompressionMap)
}

//Creates a new question and adds it to the DNS Message instance.
func (msg *Message) NewQuestion(name string, recType RecordType) {
	question := Question{}
	question.Set(name, recType)
	msg.Questions = append(msg.Questions, question)
	questionCount := msg.Header.QdCount + 1
	msg.Header.SetQuestionCount(questionCount)
}

// Checks if the message contains a question record for the given domain name.
func (msg *Message) HasQuestion(name string) bool {
	name = Canonicalize(name)
	ok := false

	if msg.Header.QdCount > 0 {
		for _, que := range msg.Questions {
			if strings.EqualFold(name, que.Name.Value) {
				ok = true
				break
			}
		} 
	}

	return ok
}

// Checks if the message contains an answer record for the given domain name.
func (msg *Message) HasAnswer(name string) bool {
	name = Canonicalize(name)
	ok := false

	if msg.Header.AnCount > 0 {
		for _, ans := range msg.Answers {
			if strings.EqualFold(name, ans.GetData()) {
				ok = true
				break
			}
		}
	}

	return ok
}

//Appends the resource records to the answers collection of the Message instance.
func (msg *Message) AddAnswers(resources []Resource) {
	msg.Answers = append(msg.Answers, resources...)
	CurrentCount := msg.Header.AnCount
	CurrentCount += uint16(len(resources))
	msg.Header.SetAnswerCount(CurrentCount)
}

//Pack the message as a sequence of octets.
func (msg *Message) Pack() []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, msg.Header.Pack()...)
	offset := len(buffer)
	if msg.Header.QdCount > 0 {
		for _, que := range msg.Questions {
			buffer = append(buffer, que.Pack(msg.compressionMap, offset)...)
			offset = len(buffer)
		}
	}

	if msg.Header.AnCount > 0 {
		for _, rr := range msg.Answers {
			buffer = append(buffer, rr.Pack(msg.compressionMap, offset)...)
			offset = len(buffer)
		}
	}

	if msg.Header.NsCount > 0 {
		for _, rr := range msg.Authoritative {
			buffer = append(buffer, rr.Pack(msg.compressionMap, offset)...)
			offset = len(buffer)
		}
	}

	if msg.Header.ArCount > 0 {
		for _, rr := range msg.Additional {
			buffer = append(buffer, rr.Pack(msg.compressionMap, offset)...)
			offset = len(buffer)
		}
	}

	return buffer
}

//Unpack the sequence of bytes to a Message instance.
func (msg *Message) Unpack(response []byte) {
	offset := 0
	offset = msg.Header.Unpack(response, offset)
	if msg.Header.QdCount > 0 {
		for index := 1; index <= int(msg.Header.QdCount); index++ {
			question := Question{}
			offset = question.Unpack(response, offset)
			msg.Questions = append(msg.Questions, question)
		}
	}

	if msg.Header.AnCount > 0 {
		for index := 1; index <= int(msg.Header.AnCount); index++ {
			answer := Resource{}
			offset = answer.Unpack(response, offset)
			msg.Answers = append(msg.Answers, answer)
		}
	}

	if msg.Header.NsCount > 0 {
		for index := 1; index <= int(msg.Header.NsCount); index++ {
			authoritative := Resource{}
			offset = authoritative.Unpack(response, offset)
			msg.Authoritative = append(msg.Authoritative, authoritative)
		}
	}

	if msg.Header.ArCount > 0 {
		for index := 1; index <= int(msg.Header.ArCount); index++ {
			additional := Resource{}
			offset = additional.Unpack(response, offset)
			msg.Additional = append(msg.Additional, additional)
		}
	}
}

//Returns a string representation of the DNS Message instance. 
func (msg *Message) String() string {
	string_value := msg.Header.String() + "\n"
	if msg.Header.QdCount > 0 {
		string_value += "QUESTION SECTION:\n"
		for _, que := range msg.Questions {
			string_value += que.String()
		}
		string_value += "\n"
	}
	
	if msg.Header.AnCount > 0 {
		string_value += "ANSWER SECTION:\n"
		for _, ans := range msg.Answers {
			string_value += ans.String()
		}
		string_value += "\n"
	}
	
	if msg.Header.NsCount > 0 {
		string_value += "AUTHORITY SECTION:\n"
		for _, auth := range msg.Authoritative {
			string_value += auth.String()
		}
		string_value += "\n"
	}

	if msg.Header.ArCount > 0 {
		string_value += "ADDITIONAL SECTION:\n"
		for _, add := range msg.Additional {
			string_value += add.String()
		}
	}

	return string_value
}

//Checks if the given resource is a response for the DNS question provided in parameter.
func (msg *Message) IsResponse(request *Message) bool {
	if !msg.Header.IsResponse {
		return false
	}

	if msg.Header.Identifier != request.Header.Identifier {
		return false
	}

	return true
}

//Returns the RRs from Answer section of DNS message matching the given record type.
func (msg *Message) FindAnswerRecords(recType RecordType) ([]Resource, bool) {
	rrValues := make([]Resource, 0)
	if msg.Header.AnCount > 0 {
		for _, ans := range msg.Answers {
			if ans.Type == recType {
				rrValues = append(rrValues, ans)
			}
		}
	}

	if len(rrValues) > 0 {
		return rrValues, true
	} else {
		return nil, false
	}
}

//Returns the RRs from Authoritative section of DNS message matching the given record type.
func (msg *Message) FindAuthorityRecords(recType RecordType) ([]Resource, bool) {
	rrValues := make([]Resource, 0)
	if msg.Header.NsCount > 0 {
		for _, ns := range msg.Authoritative {
			if ns.Type == recType {
				rrValues = append(rrValues, ns)
			}
		}
	}

	if len(rrValues) > 0 {
		return rrValues, true
	} else {
		return nil, false
	}
}

//Returns the RRs from Additional section of DNS message matching the given record type.
func (msg *Message) FindAdditionalRecords(recType RecordType) ([]Resource, bool) {
	rrValues := make([]Resource, 0)
	if msg.Header.ArCount > 0 {
		for _, ans := range msg.Additional {
			if ans.Type == recType {
				rrValues = append(rrValues, ans)
			}
		}
	}

	if len(rrValues) > 0 {
		return rrValues, true
	} else {
		return nil, false
	}
}