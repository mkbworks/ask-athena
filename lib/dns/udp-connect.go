package dns

import (
	"net"
	"strconv"
	"errors"
)

//Structure to manage a single UDP connection.
type UdpConnect struct {
	Connection *net.UDPConn
}

//Uses User Datagram Protocol (UDP) to connect to the remote server and port number and return the connection object.
func (uc *UdpConnect) ConnectTo(RemoteAddress string, PortNumber int) {
	address_string := RemoteAddress + ":" + strconv.Itoa(PortNumber)
	udpAddr, err := net.ResolveUDPAddr(MESSAGE_PROTOCOL, address_string)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP(MESSAGE_PROTOCOL, nil, udpAddr)
	if err != nil {
		panic(err)
	}

	uc.Connection = conn
}

//Sends the given byte stream across the UDP connection.
func (uc *UdpConnect) Send(buffer []byte) {
	if len(buffer) > UDP_MESSAGE_SIZE_LIMIT {
		panic(errors.New("request message size exceeds 512 bytes"))
	}
	
	_, err := uc.Connection.Write(buffer)
	if err != nil {
		panic(err)
	}
}

//Receives a stream of bytes from the UDP connection.
func (uc *UdpConnect) Receive() []byte {
	buffer := make([]byte, UDP_MESSAGE_SIZE_LIMIT)
	byteCount, err := uc.Connection.Read(buffer)
	if err != nil {
		panic(err)
	}
	return buffer[:byteCount]
}

//Close the given UDP connection.
func (uc *UdpConnect) Close() {
	err := uc.Connection.Close()
	if err != nil {
		panic(err)
	}
}