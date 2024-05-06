package dns

import (
	"errors"
	"fmt"
	"net"
)

//Structure to represent a DNS Resolver.
type Resolver struct {
	RemoteServer *net.UDPConn
	AllowedRRTypes RecordTypes
}

//Queries the DNS server and fetches the 't' type record for hostname - 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) []string {
	resolvedValue := make([]string, 0)
	recordType := resolver.AllowedRRTypes.GetKey(t)
	if recordType != "" {
		fmt.Printf("Attempting to fetch '%s' type record for %s", recordType, name)
		request := Message{}
		response := Message{}
		request.SetRequest(name, t)
		requestBuffer := request.Pack()
		fmt.Println("The request buffer generated is:", requestBuffer)
		resolver.Send(requestBuffer)
		responseBuffer := resolver.Receive()
		response.Unpack(responseBuffer)
	} else {
		panic(errors.New("given record type is not one of the acceptable types"))
	}

	return resolvedValue
}

//Sends the request stream generated to the DNS server.
func (resolver *Resolver) Send(request []byte) {
	if len(request) > UDP_MESSAGE_SIZE_LIMIT {
		panic(errors.New("request message size exceeds 512 bytes"))
	}
	
	_, err := resolver.RemoteServer.Write(request)
	if err != nil {
		panic(err)
	}
}

//Receives the response back from the DNS Server.
func (resolver *Resolver) Receive() []byte {
	buffer := make([]byte, UDP_MESSAGE_SIZE_LIMIT)
	byteCount, err := resolver.RemoteServer.Read(buffer)
	if err != nil {
		panic(err)
	}
	return buffer[:byteCount]
}


//Closes the RemoteServer instance associated with the resolver.
func (resolver *Resolver) Close() {
	resolver.RemoteServer.Close()
}