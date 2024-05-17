package dns

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//Returns a new instance of Resolver. In case of any errors, it returns nil instead.
func GetResolver(RootServersPath string, CacheFilePath string, LogFilePath string) (*Resolver, error) {
	isRootServerAbs := filepath.IsAbs(RootServersPath)
	isCacheFilePathAbs := filepath.IsAbs(CacheFilePath)
	if !isRootServerAbs {
		return nil, ErrNotAbsolutePath
	}

	if !isCacheFilePathAbs {
		return nil, ErrNotAbsolutePath
	}
	resolver := Resolver{}
	resolver.RootServers = BindFile{}
	err := resolver.RootServers.Initialize(RootServersPath)
	if err != nil {
		return nil, err
	}
	resolver.Cache = BindFile{}
	err = resolver.Cache.Initialize(CacheFilePath)
	if err != nil {
		return nil, err
	}
	logFileHandler, err := os.Create(LogFilePath)
	if err != nil {
		return nil, err
	}
	resolver.Logger = log.New(logFileHandler, "", log.Ldate | log.Ltime)
	resolver.Logger.Println("Local records have been moved from file to memory.")
	resolver.Logger.Println("Root DNS Server records have been moved from BIND file to memory.")
	return &resolver, nil
}

//Generates a random 16-bit integer as DNS Message Id.
func Id() uint16 {
	var Identifier uint16
	binary.Read(rand.Reader, binary.BigEndian, &Identifier)
	return Identifier
}

//Packs an 16-bit unsigned integer into a stream of octets (or bytes) and returns them as an array of byte values.
func PackUInt16(number uint16) []byte {
	binary_value := fmt.Sprintf("%b", number)
	return PackBinary16(binary_value)
}

//Packs a 16-bit binary value string into a stream of octets (aka bytes) and returns them as an array of byte values.
func PackBinary16(binary_value string) []byte {
	buffer := make([]byte, 0)
	if len(binary_value) < 16 {
		padding_length := 16 - len(binary_value)
		binary_value = strings.Repeat("0", padding_length) + binary_value
	}
	ms_octet_string := binary_value[:8]
	ls_octet_string := binary_value[8:]
	ms_octet_value, _ := strconv.ParseUint(ms_octet_string, 2, 8)
	ls_octet_value, _ := strconv.ParseUint(ls_octet_string, 2, 8)
	buffer = append(buffer, byte(ms_octet_value))
	buffer = append(buffer, byte(ls_octet_value))
	return buffer
}

//Unpacks a stream of bytes into a uint16 number.
func UnpackUInt16(buffer []byte) uint16 {
	return_value := UnpackBinary(buffer)
	number_value, _ := strconv.ParseUint(return_value, 2, 16)
	return uint16(number_value)
}

//Unpacks a stream of bytes into a uint32 number.
func UnpackUInt32(buffer []byte) uint32 {
	return_value := UnpackBinary(buffer)
	number_value, _ := strconv.ParseUint(return_value, 2, 32)
	return uint32(number_value)
}

//Unpacks a stream of bytes into a binary string value.
func UnpackBinary(buffer []byte) string {
	return_value := ""
	for _, octet := range buffer {
		binary_string := fmt.Sprintf("%b", octet)
		if len(binary_string) < 8 {
			padding_length := 8 - len(binary_string)
			binary_string = strings.Repeat("0", padding_length) + binary_string
		}
		return_value += binary_string
	}
	return return_value
}

//Returns the 'bit_count' least significant bits from the binary representation of 'number'.
func GetBinary(number uint16, bit_count int) string {
	binary_value := fmt.Sprintf("%b", number)
	if len(binary_value) < bit_count {
		padding_length := bit_count - len(binary_value)
		binary_value = strings.Repeat("0", padding_length) + binary_value
	} else if len(binary_value) > bit_count {
		return ""
	}
	return binary_value
}

//Creates and returns a new Message instance.
func NewMessage(mt MessageType) *Message {
	message := Message{}
	message.Initialize(mt)
	return &message
}

//Creates a new resource record and returns a pointer to the Resource instance.
func NewResourceRecord(dname string, ttl uint32, class string, recType string, data string) *Resource {
	resource := Resource{}
	resource.Initialize(dname, recType, class, ttl, data)
	return &resource
}

//Parses the given string and returns its uint64 equivalent.
func parseUIntString(value string, bitsize int) uint64 {
	conv_value, _ := strconv.ParseUint(value, 10, bitsize)
	return conv_value
}

//Converts the given byte stream to an IP Address string
func getIPAddress(buffer []byte) string {
	IP := net.IP(buffer)
	return IP.String()
}

//Converts the given IP address string to a stream of bytes.
func convertToBytes(IpAddress string, IpType string) []byte {
	emptyResponse := make([]byte, 0)
	ip := net.ParseIP(IpAddress)
	if ip == nil {
		return emptyResponse
	}
	
	if strings.EqualFold(IpType, "IPv4") {
		ipv4 := ip.To4()
		if ipv4 == nil {
			return emptyResponse
		}
		return ipv4
	} else if strings.EqualFold(IpType, "IPv6") {
		ipv6 := ip.To16()
		if ipv6 == nil {
			return emptyResponse
		}
		return ipv6
	} else {
		return emptyResponse
	}
}

//Returns a string representing the canonicalized value of given domain name.
func Canonicalize(domainName string) string {
	domainName = strings.Trim(domainName, DOMAIN_LABEL_SEPERATOR)
	domainName = strings.ToLower(domainName)
	domainName += DOMAIN_LABEL_SEPERATOR
	return domainName
}