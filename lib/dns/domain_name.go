package dns

import (
	"strings"
)

//Represents a domain name as defined in RFC 1035.
type DomainName struct {
	//Contains a byte stream representation of the domain name.
	Data []byte
	//Contains the string value representation of the domain name
	Value string
	//Contains the number of labels (domains & subdomains) in the domain name.
	Length uint8
}

//Initialises the instance of DomainName
func (name *DomainName) Initialize(dName string) {
	name.Data = make([]byte, 0)
	dName = strings.Trim(dName, DOMAIN_LABEL_SEPERATOR)
	name.Length = uint8(len(strings.Split(dName, DOMAIN_LABEL_SEPERATOR)))
	dName = Canonicalize(dName)
	name.Value = dName
}

//Parse the given domain name string and pack it as sequence of octets.
func (name *DomainName) Pack(compressionMap CompressionMap, offset int) []byte {
	encodedBytes := make([]byte, 0)
	dName := name.Value
	isPtrAvailable := false

	for {
		if dName == "" {
			encodedBytes = append(encodedBytes, byte(0))
			break
		}
		new_offset, ok := compressionMap[dName]
		if ok {
			ptrUIntValue := PTR_DETECT_VALUE | uint16(new_offset)
			encodedBytes = append(encodedBytes, PackUInt16(ptrUIntValue)...)
			isPtrAvailable = true
			break
		}

		label, pendingDName, ok := strings.Cut(dName, DOMAIN_LABEL_SEPERATOR)
		dName = pendingDName
		if ok {
			labelBytes := []byte(label)
			encodedBytes = append(encodedBytes, byte(len(labelBytes)))
			encodedBytes = append(encodedBytes, labelBytes...)
		}
	}

	if !isPtrAvailable {
		compressionMap[dName] = offset
	}

	name.Data = encodedBytes
	return encodedBytes
}

//Gets the byte length of the given domain name.
func (name *DomainName) GetLength() int {
	length := 0
	dName := name.Value
	for _, dLabel := range strings.Split(dName, DOMAIN_LABEL_SEPERATOR) {
		length +=1 //for the byte representing the length of the label or subdomain
		length += len([]byte(dLabel))
	}
	length +=1 //For the final 0 byte present to signify end of domain name
	return length
}

//Unpack the given byte stream and extract the domain name.
func (name *DomainName) Unpack(buffer []byte, offset int) int {
	completeDomainName, offset := name.getDomainName(buffer, offset)
	name.Value = completeDomainName
	name.Length = uint8(len(strings.Split(name.Value, DOMAIN_LABEL_SEPERATOR)))
	return offset
}

//Parses the given byte stream and fetches the domain name. Domain name can be represented directly
//or can be compressed and represented through a pointer as per RFC 1035 - Section 4.1.4
func (name *DomainName) getDomainName(buffer []byte, offset int) (string, int) {
	completeDomainName := ""
	LastIndexRead := offset
	labelByteCount := buffer[LastIndexRead]
	for iterate := true; iterate; {
		PtrBytesCheck := buffer[LastIndexRead: LastIndexRead + 2]
		PtrBytesValue := UnpackUInt16(PtrBytesCheck)
		if PtrBytesValue & PTR_DETECT_VALUE == PTR_DETECT_VALUE {
			ptr_offset_value := PtrBytesValue & PTR_OFFSET_FETCH
			subdomain, _ := name.getDomainName(buffer, int(ptr_offset_value))
			completeDomainName = completeDomainName + DOMAIN_LABEL_SEPERATOR + subdomain
			LastIndexRead = LastIndexRead + 2
			iterate = false
		} else if int(labelByteCount) != 0 {
			labelBytes := buffer[LastIndexRead + 1: LastIndexRead + int(labelByteCount) + 1]
			name.Data = append(name.Data, labelByteCount)
			name.Data = append(name.Data, labelBytes...)
			completeDomainName = completeDomainName + DOMAIN_LABEL_SEPERATOR + string(labelBytes)
			LastIndexRead = LastIndexRead + int(labelByteCount) + 1
			labelByteCount = buffer[LastIndexRead]
		} else {
			name.Data = append(name.Data, byte(0))
			LastIndexRead = LastIndexRead + 1
			iterate = false
		}
	}

	completeDomainName = Canonicalize(completeDomainName)
	return completeDomainName, LastIndexRead
}

//Returns the domain name as a string.
func (name *DomainName) String() string {
	return name.Value
}