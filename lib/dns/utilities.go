package dns

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//Returns a new instance of Resolver. In case of any errors, it returns nil instead.
func NewResolver(RootServersPath string, CacheFilePath string, traceLogs bool) (*Resolver, error) {
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
	resolver.Logger = log.New(os.Stdout, "", log.Ldate | log.Ltime)
	resolver.traceLogs = traceLogs
	resolver.response = nil
	return &resolver, nil
}

//Generates a random 16-bit integer as DNS Message Id.
func Id() uint16 {
	var Identifier uint16
	binary.Read(rand.Reader, binary.BigEndian, &Identifier)
	return Identifier
}

//Packs a 16-bit unsigned integer into a stream of octets (or bytes) and returns them as an array of byte values.
func PackUInt32(number uint32) []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, number)
	return buffer
}

//Packs a 16-bit unsigned integer into a stream of octets (or bytes) and returns them as an array of byte values.
func PackUInt16(number uint16) []byte {
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, number)
	return buffer
}

//Unpacks a stream of bytes into a uint16 number.
func UnpackUInt16(buffer []byte) uint16 {
	return binary.BigEndian.Uint16(buffer)
}

//Unpacks a stream of bytes into a uint32 number.
func UnpackUInt32(buffer []byte) uint32 {
	return binary.BigEndian.Uint32(buffer)
}

//Creates and returns a new Message instance.
func NewMessage(mt MessageType, MsgId uint16) *Message {
	message := Message{}
	message.Initialize(mt, MsgId)
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
	
	if strings.EqualFold(IpType, ADDRESS_IPv4) {
		ipv4 := ip.To4()
		if ipv4 == nil {
			return emptyResponse
		}
		return ipv4
	} else if strings.EqualFold(IpType, ADDRESS_IPv6) {
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