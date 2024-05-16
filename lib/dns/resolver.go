package dns

import (
	"errors"
)

//Type to denote the type of DNS Server to send request to.
type ServerType uint8

//Structure to represent a DNS Resolver.
type Resolver struct {
	//References the BIND file containing the DNS root server details.
	RootServers BindFile
	//References the BIND file containing all the cached resource records.
	Cache BindFile
}

//Queries the DNS server and fetches the 't' type record for 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) []string {
	recordType := AllowedRRTypes.GetKey(t)
	if recordType != "" {
		rootServerAddresses := resolver.getRootServer(TYPE_A)		
		TldNameServers := resolver.GetNameServer(name, t, rootServerAddresses)
		AuthNameServers := resolver.GetNameServer(name, t, TldNameServers)
		values := resolver.GetAuthoritativeData(name, t, AuthNameServers)
		return values
	} else {
		panic(errors.New("given record type is not one of the acceptable types"))
	}
}

//Returns true if the record type provided is accepted by the resolver, else returns false.
func (resolver *Resolver) IsAllowed(recordType string) bool {
	_, exists := AllowedRRTypes[recordType]
	return exists
}

//Returns the record type object for the given type string.
func (resolver *Resolver) GetRecordType(recordType string) RecordType {
	return AllowedRRTypes.GetRecordType(recordType)
}

//Returns the list of all IPv4 addresses for the root DNS server.
func (resolver *Resolver) getRootServer(recType RecordType) []string {
	rootServerAddress := make([]string, 0)
	for _, rr := range resolver.RootServers.ResourceRecords {
		if rr.Type == recType {
			rootServerAddress = append(rootServerAddress, rr.GetData())
		}
	}

	if len(rootServerAddress) == 0 {
		panic(errors.New("root dns server - ipv4 address not found"))
	}

	return rootServerAddress
}

//Returns the name server details returned by either the Root DNS or TLD Name servers for the domain name.
func (resolver *Resolver) GetNameServer(name string, t RecordType, servers []string) []string {
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, t)
	NameServers := make([]string, 0)
	for _, address := range servers {
		response := resolver.GetResponse(request, address)
		NsHosts, isExists := response.GetRRValuesFor(name, TYPE_NS)
		if !isExists {
			continue
		}

		for _, NsHost := range NsHosts {
			NsIps, exists := response.GetRRValuesFor(NsHost, TYPE_A)
			if exists {
				NameServers = append(NameServers, NsIps...)
			} else {
				values := resolver.ResolveNameServer(NsHost, servers)
				NameServers = append(NameServers, values...)
			}
		}

		if len(NameServers) > 0 {
			break
		}
	}

	return NameServers
}

//Returns the requested information for the given domain name by querying the authoritative DNS servers.
func (resolver *Resolver) GetAuthoritativeData(name string, t RecordType, servers []string) []string {
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, t)
	AuthoritativeData := make([]string, 0)
	for _, address := range servers {
		response := resolver.GetResponse(request, address)
		values, exists := response.GetRRValuesFor(name, t)
		if exists {
			return values
		}

		NsHosts, exists := response.GetRRValuesFor(name, TYPE_NS)
		if exists {
			NameServerIps := make([]string, 0)
			for _, NsHost := range NsHosts {
				NsIps, IsAlreadyThere := response.GetRRValuesFor(NsHost, TYPE_A)
				if IsAlreadyThere {
					NameServerIps = append(NameServerIps, NsIps...)
				} else {
					NameServerIps = append(NameServerIps, resolver.ResolveNameServer(NsHost, servers)...)
				}

				if len(NameServerIps) > 0 {
					break
				}
			}

			return resolver.GetAuthoritativeData(name, t, NameServerIps) 
		}

		Cnames, exists := response.GetRRValuesFor(name, TYPE_CNAME)
		if exists {
			for _, cname := range Cnames {
				resp := response.GetIPForCNAME(cname, t)
				if len(resp) > 0 {
					AuthoritativeData = append(AuthoritativeData, resp...)
				} else {
					AuthoritativeData = append(AuthoritativeData, resolver.GetAuthoritativeData(cname, t, servers)...)
				}
			}
		}

		if len(AuthoritativeData) > 0 {
			return AuthoritativeData
		}
	}

	return AuthoritativeData
}

//Resolves the given name server and returns its IPv4 addresses. This function is written with the assumption that the name server host name itself does not have a CNAME record attached to it and can directly fetch type-A record.
func (resolver *Resolver) ResolveNameServer(name string, servers []string) []string {
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, TYPE_A)
	IPv4Addresses := make([]string, 0)
	for _, server := range servers {
		response := resolver.GetResponse(request, server)
		values, exists := response.GetRRValuesFor(name, TYPE_A)
		if exists {
			IPv4Addresses = append(IPv4Addresses, values...)
			break
		}
		/*
		values, exists = response.GetRRValuesFor(name, TYPE_NS)
		if exists {
			for _, value := range values {
				IPv4Addresses = append(IPv4Addresses, resolver.ResolveNameServer(value, servers)...)
			}
		}
		if len(IPv4Addresses) > 0 {
			break
		} */
	}

	return IPv4Addresses
}

//Sends the request to the target DNS server and receives a response over the same connection.
func (resolver *Resolver) GetResponse(request *Message, ServerAddress string) *Message {
	SendBuffer := request.Pack()
	udpConnect := UdpConnect{}
	udpConnect.ConnectTo(ServerAddress, DNS_PORT_NUMBER)
	var response *Message
	for validResponse := false; !validResponse; {
		udpConnect.Send(SendBuffer)
		receiveBuffer := udpConnect.Receive()
		response = NewMessage(MSG_RESPONSE)
		response.Unpack(receiveBuffer)
		validResponse = response.IsResponse(request)
	}
	udpConnect.Close()
	return response
}