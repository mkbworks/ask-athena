package dns

import (
	"net"
	"strconv"
)

//Structure to manage a single UDP connection.
type UdpConnect struct {
	Connection *net.UDPConn
}

//Uses User Datagram Protocol (UDP) to connect to the remote server and port number and return the connection object.
func (uc *UdpConnect) ConnectTo(RemoteAddress string, PortNumber int) error {
	address_string := RemoteAddress + ":" + strconv.Itoa(PortNumber)
	udpAddr, err := net.ResolveUDPAddr(MESSAGE_PROTOCOL, address_string)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP(MESSAGE_PROTOCOL, nil, udpAddr)
	if err != nil {
		return err
	}

	uc.Connection = conn
	return nil
}

//Sends the given byte stream across the UDP connection.
func (uc *UdpConnect) Send(buffer []byte) error {
	if len(buffer) > UDP_MESSAGE_SIZE_LIMIT {
		return ErrMessageTooLong
	}
	
	_, err := uc.Connection.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

//Receives a stream of bytes from the UDP connection.
func (uc *UdpConnect) Receive() ([]byte, error) {
	buffer := make([]byte, UDP_MESSAGE_SIZE_LIMIT)
	byteCount, err := uc.Connection.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:byteCount], nil
}

//Close the given UDP connection.
func (uc *UdpConnect) Close() error {
	err := uc.Connection.Close()
	if err != nil {
		return err
	}

	return nil
}