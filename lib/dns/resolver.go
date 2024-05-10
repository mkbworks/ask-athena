package dns

import (
	"errors"
	"net"
	"fmt"
)

//Structure to represent a DNS Resolver.
type Resolver struct {
	RemoteServer *net.UDPConn
	AllowedRRTypes RecordTypes
}

//Queries the DNS server and fetches the 't' type record for hostname - 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) {
	recordType := resolver.AllowedRRTypes.GetKey(t)
	if recordType != "" {
		fmt.Printf("Attempting to fetch '%s' type record for %s", recordType, name)
		request := GetMessage(MSG_REQUEST)
		fmt.Println("New DNS Request message has been created.")
		response := GetMessage(MSG_RESPONSE)
		request.NewQuestion(name, t)
		fmt.Println("New question has been created and added to the DNS Request.")
		requestBuffer := request.Pack()
		fmt.Println("The request buffer generated is:", requestBuffer)
		resolver.Send(requestBuffer)
		fmt.Println("DNS Request message has been sent to the target UDP Server.")
		responseBuffer := resolver.Receive()
		fmt.Println("The response buffer received is:", responseBuffer)
		response.Unpack(responseBuffer)
		fmt.Println("Answer received.")
		fmt.Println(response.String())
	} else {
		panic(errors.New("given record type is not one of the acceptable types"))
	}
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

//Returns true if the record type provided is accepted by the resolver, else returns false.
func (resolver *Resolver) IsAllowed(recordType string) bool {
	_, exists := resolver.AllowedRRTypes[recordType]
	return exists
}

//Returns the record type object for the given type string.
func (resolver *Resolver) GetRecordType(recordType string) RecordType {
	rt, exists := resolver.AllowedRRTypes[recordType]
	if exists {
		return rt
	} else {
		panic(errors.New("record type could not be found in the list of record types supported by DNS resolver"))
	}
}

//Closes the RemoteServer instance associated with the resolver.
func (resolver *Resolver) Close() {
	resolver.RemoteServer.Close()
}