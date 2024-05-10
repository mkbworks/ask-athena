package dns

import (
	"strconv"
	"strings"
	"fmt"
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
func (name *DomainName) Initialize() {
	name.Data = make([]byte, 0)
}

//Parse the given domain name string and pack it as sequence of octets.
func (name *DomainName) Pack(dName string) {
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

//Unpack the given byte stream and extract the domain name.
func (name *DomainName) Unpack(buffer []byte, offset int) int {
	completeDomainName, offset := name.getDomainName(buffer, offset)
	name.Value = completeDomainName
	name.Length = uint8(len(strings.Split(name.Value, DOMAIN_LABEL_SEPERATOR)))
	return offset
}

//Parses the given byte stream and fetches the domain name string. Domain name can be represented directly
//or can be compressed and represented through a pointer as per RFC 1035 - Section 4.1.4
func (name *DomainName) getDomainName(buffer []byte, offset int) (string, int) {
	fmt.Printf("Starting to unpack Domain name from offset = %d\n", offset)
	completeDomainName := ""
	fmt.Println("Complete domain name at the start:", completeDomainName)
	LastIndexRead := offset
	labelByteCount := buffer[LastIndexRead]
	for iterate := true; iterate; {
		PtrBytesCheck := buffer[LastIndexRead: LastIndexRead + 2]
		binary_string := UnpackBinary(PtrBytesCheck)
		if strings.HasPrefix(binary_string, DOMAIN_NAME_PTR_PREFIX) {
			ptr_offset_bits := binary_string[len(DOMAIN_NAME_PTR_PREFIX):]
			ptr_offset_value, err := strconv.ParseUint(ptr_offset_bits, 2, 16)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Pointer bits found. Pointer offset value is %d\n", int(ptr_offset_value))
			subdomain, _ := name.getDomainName(buffer, int(ptr_offset_value))
			completeDomainName = completeDomainName + DOMAIN_LABEL_SEPERATOR + subdomain
			LastIndexRead = LastIndexRead + 2
			iterate = false
		} else if int(labelByteCount) != 0 {
			fmt.Printf("Parsing subdomain from offset = %d\n", LastIndexRead + 1)
			labelBytes := buffer[LastIndexRead + 1: LastIndexRead + int(labelByteCount) + 1]
			name.Data = append(name.Data, labelByteCount)
			name.Data = append(name.Data, labelBytes...)
			completeDomainName = completeDomainName + DOMAIN_LABEL_SEPERATOR + string(labelBytes)
			LastIndexRead = LastIndexRead + int(labelByteCount) + 1
			labelByteCount = buffer[LastIndexRead]
		} else {
			LastIndexRead = LastIndexRead + 1
			iterate = false
		}
		fmt.Println("Complete domain name at the end of an iteration:", completeDomainName)
	}

	completeDomainName = strings.TrimPrefix(strings.TrimSuffix(completeDomainName, DOMAIN_LABEL_SEPERATOR), DOMAIN_LABEL_SEPERATOR)
	return completeDomainName, LastIndexRead
}

//Implements the GoStringer interface to provide a string implementation of domain name.
func (name *DomainName) String() string {
	return name.Value
}