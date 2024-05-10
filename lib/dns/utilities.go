package dns

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

//Returns a new instance of Resolver. In case of any errors, it returns nil instead.
func GetResolver() (*Resolver, error) {
	resolver := Resolver{}
	resolver.AllowedRRTypes = RecordTypes{
		"A":     TYPE_A,
		"NS":    TYPE_NS,
		"CNAME": TYPE_CNAME,
		"TXT":   TYPE_TXT,
		"AAAA":  TYPE_AAAA,
	}
	address_string := ROOT_SERVER_ADDRESS + ":" + strconv.Itoa(DNS_PORT_NUMBER)
	udpAddr, err := net.ResolveUDPAddr(MESSAGE_PROTOCOL, address_string)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(MESSAGE_PROTOCOL, nil, udpAddr)
	if err != nil {
		return nil, err
	}
	resolver.RemoteServer = conn
	return &resolver, nil
}

//Generates a random 16-bit integer as DNS Message Id.
func Id() uint16 {
	var Identifier uint16
	err := binary.Read(rand.Reader, binary.BigEndian, &Identifier)
	if err != nil {
		panic(err)
	}
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
	ms_octet_value, err := strconv.ParseUint(ms_octet_string, 2, 8)
	if err != nil {
		panic(err)
	}
	ls_octet_value, err := strconv.ParseUint(ls_octet_string, 2, 8)
	if err != nil {
		panic(err)
	}
	buffer = append(buffer, byte(ms_octet_value))
	buffer = append(buffer, byte(ls_octet_value))
	return buffer
}

//Unpacks a stream of bytes into a uint16 number.
func UnpackUInt16(buffer []byte) uint16 {
	return_value := UnpackBinary(buffer)
	number_value, err := strconv.ParseUint(return_value, 2, 16)
	if err != nil {
		panic(err)
	}
	return uint16(number_value)
}

//Unpacks a stream of bytes into a uint32 number.
func UnpackUInt32(buffer []byte) uint32 {
	return_value := UnpackBinary(buffer)
	number_value, err := strconv.ParseUint(return_value, 2, 32)
	if err != nil {
		panic(err)
	}
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
		panic(errors.New("bit count for the given number is larger than the required bit count"))
	}
	return binary_value
}

//Creates and returns a new Message instance.
func GetMessage(mt MessageType) *Message {
	message := Message{}
	message.Initialize(mt)
	return &message
}