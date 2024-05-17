package dns

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
)

//Type to denote the type of DNS Server to send request to.
type ServerType uint8

//Structure to represent a DNS Resolver.
type Resolver struct {
	//References the BIND file containing the DNS root server details.
	RootServers BindFile
	//References the BIND file containing all the cached resource records.
	Cache BindFile
	//Logger to be used to generate logs.
	Logger *log.Logger
}

//Queries the DNS server and fetches the 't' type record for 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) []string {
	emptyIPs := make([]string, 0)
	recordType := AllowedRRTypes.GetKey(t)
	if recordType != "" {
		if t == TYPE_A {
			cacheRecords, ok := resolver.Cache.FindAll(name, TYPE_A)
			if ok {
				resolver.Logger.Println("IPv4 address for the given domain name has been served from local Cache file.")
				return cacheRecords
			}
			A_RRs, err := resolver.resolveA(name)
			if err != nil {
				resolver.Logger.Println(err.Error())
				return make([]string, 0)
			}

			values := make([]string, 0)
			for _, A_RR := range A_RRs {
				values = append(values, A_RR.GetData())
				resolver.Cache.Add(A_RR.Name.Value, A_RR.TTL, A_RR.Class.String(), A_RR.Type.String(), A_RR.GetData())
			}

			return values
		} else if t == TYPE_AAAA {
			cacheRecords, ok := resolver.Cache.FindAll(name, TYPE_AAAA)
			if ok {
				resolver.Logger.Println("IPv6 address for the given domain name has been served from local Cache file.")
				return cacheRecords
			}
			AAAA_RRs, err := resolver.resolveAAAA(name)
			if err != nil {
				resolver.Logger.Println(err.Error())
				return make([]string, 0)
			}

			values := make([]string, 0)
			for _, AAAA_RR := range AAAA_RRs {
				values = append(values, AAAA_RR.GetData())
				resolver.Cache.Add(AAAA_RR.Name.Value, AAAA_RR.TTL, AAAA_RR.Class.String(), AAAA_RR.Type.String(), AAAA_RR.GetData())
			}

			return values
		} else if t == TYPE_TXT {
			value := resolver.resolveTXT(name)
			values := make([]string, 0)
			values = append(values, fmt.Sprintf("%d", value))
			return values
		} else if t == TYPE_CNAME {
			cacheRecords, ok := resolver.Cache.FindAll(name, TYPE_CNAME)
			if ok {
				values := make([]string, 0)
				CNAME_IPs, ok := resolver.Cache.FindAll(cacheRecords[0], TYPE_A)
				if ok {
					values = append(values, cacheRecords[0])
					values = append(values, CNAME_IPs...)
					resolver.Logger.Println("CNAME and its IPv4 address have been served from local Cache file.")
					return values
				}
			}
			CNAME_RRs, err := resolver.resolveCNAME(name)
			if err != nil {
				resolver.Logger.Println(err.Error())
				return make([]string, 0)
			}

			if len(CNAME_RRs) > 0 {
				for _, CNAME_RR := range CNAME_RRs {
					resolver.Cache.Add(CNAME_RR.Name.Value, CNAME_RR.TTL, CNAME_RR.Class.String(), CNAME_RR.Type.String(), CNAME_RR.GetData())
				}
				values := make([]string, 0)
				values = append(values, CNAME_RRs[0].GetData())
				CNAME_IPs, err := resolver.resolveA(CNAME_RRs[0].GetData())
				if err != nil {
					resolver.Logger.Println(err.Error())
					return make([]string, 0)
				}
				for _, CNAME_IP := range CNAME_IPs {
					values = append(values, CNAME_IP.GetData())
					resolver.Cache.Add(CNAME_IP.Name.Value, CNAME_IP.TTL, CNAME_IP.Class.String(), CNAME_IP.Type.String(), CNAME_IP.GetData())
				}
				return values
			} else {
				return make([]string, 0)
			}
		}
		
	} else {
		resolver.Logger.Println(ErrInvalidRecordType.Error())
	}

	return emptyIPs
}

func (resolver *Resolver) resolveA(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	for {
		request := NewMessage(MSG_REQUEST)
		request.NewQuestion(name, TYPE_A)
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				return resolver.resolveA(CNAME_RRs[0].GetData())
			}

			A_RRs, _ := response.FindAnswerRecords(TYPE_A)
			return A_RRs, nil
		}

		if response.Header.NsCount > 0 && response.Header.ArCount == 0 {
			NS_RRs, Exists := response.FindAuthorityRecords(TYPE_NS)
			if !Exists {
				return nil, ErrNameServerFetch
			}
			NS_IPs, err := resolver.resolveA(NS_RRs[0].GetData())
			if err != nil {
				return nil, err
			}

			nameserver = NS_IPs[0].GetData()
		} else {
			Add_RRs, Exists := response.FindAdditionalRecords(TYPE_A)
			if !Exists {
				return nil, ErrNameServerFetch
			}

			nameserver = Add_RRs[0].GetData()
		}
	}
}

func (resolver *Resolver) resolveAAAA(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	for {
		request := NewMessage(MSG_REQUEST)
		request.NewQuestion(name, TYPE_AAAA)
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				return resolver.resolveAAAA(CNAME_RRs[0].GetData())
			}

			AAAA_RRs, _ := response.FindAnswerRecords(TYPE_AAAA)
			return AAAA_RRs, nil
		}

		if response.Header.NsCount > 0 && response.Header.ArCount == 0 {
			NS_RRs, Exists := response.FindAuthorityRecords(TYPE_NS)
			if !Exists {
				return nil, ErrNameServerFetch
			}
			NS_IPs, err := resolver.resolveA(NS_RRs[0].GetData())
			if err != nil {
				return nil, err
			}

			nameserver = NS_IPs[0].GetData()
		} else {
			Add_RRs, Exists := response.FindAdditionalRecords(TYPE_A)
			if !Exists {
				return nil, ErrNameServerFetch
			}

			nameserver = Add_RRs[0].GetData()
		}
	}
}

func (resolver *Resolver) resolveTXT(name string) int {
	RandomNumber := rand.IntN(10000 * len(name))
	if RandomNumber < 1000 {
		RandomNumber += 1000
	}
	return RandomNumber
}

func (resolver *Resolver) resolveCNAME(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	for {
		request := NewMessage(MSG_REQUEST)
		request.NewQuestion(name, TYPE_CNAME)
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, Exists := response.FindAnswerRecords(TYPE_CNAME)
			if Exists {
				return CNAME_RRs, nil
			} else {
				return make([]Resource, 0), nil
			}
		}

		if response.Header.NsCount > 0 && response.Header.ArCount == 0 {
			NS_RRs, Exists := response.FindAuthorityRecords(TYPE_NS)
			if !Exists {
				return nil, ErrNameServerFetch
			}
			NS_IPs, err := resolver.resolveA(NS_RRs[0].GetData())
			if err != nil {
				return nil, err
			}

			nameserver = NS_IPs[0].GetData()
		} else {
			Add_RRs, Exists := response.FindAdditionalRecords(TYPE_A)
			if !Exists {
				return nil, ErrNameServerFetch
			}

			nameserver = Add_RRs[0].GetData()
		}
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

//Returns the IP address of a root DNS server chosen at random.
func (resolver *Resolver) getRootServer(recType RecordType) string {
	rootServerAddress := make([]string, 0)
	for _, rr := range resolver.RootServers.ResourceRecords {
		if rr.resource.Type == recType {
			rootServerAddress = append(rootServerAddress, rr.resource.GetData())
		}
	}

	RandomIndex := rand.IntN(len(rootServerAddress))
	RandomServer := rootServerAddress[RandomIndex]

	return RandomServer
}


//Sends the request to the target DNS server and receives a response over the same connection.
func (resolver *Resolver) getResponse(request *Message, ServerAddress string) *Message {
	ServerAddress = strings.TrimSpace(ServerAddress)
	if ServerAddress == "" {
		return nil
	}
	resolver.Logger.Printf("**********************************************\n")
	resolver.Logger.Printf("DNS Request being sent to server - %s.\n", ServerAddress)
	resolver.Logger.Printf("**********************************************\n")
	resolver.Logger.Printf("Request Contents are:\n%s", request.String())
	resolver.Logger.Printf("**********************************************\n")
	SendBuffer := request.Pack()
	udpConnect := UdpConnect{}
	err := udpConnect.ConnectTo(ServerAddress, DNS_PORT_NUMBER)
	if err != nil {
		resolver.Logger.Println(err.Error())
		return nil
	}
	var response *Message
	for validResponse := false; !validResponse; {
		err = udpConnect.Send(SendBuffer)
		if err != nil {
			resolver.Logger.Println(err.Error())
			return nil
		}
		receiveBuffer, err := udpConnect.Receive()
		if err != nil {
			resolver.Logger.Println(err.Error())
			return nil
		}
		response = NewMessage(MSG_RESPONSE)
		response.Unpack(receiveBuffer)
		validResponse = response.IsResponse(request)
	}
	resolver.Logger.Printf("Response received back:\n%s", response.String())
	resolver.Logger.Printf("**********************************************\n")
	err = udpConnect.Close()
	if err != nil {
		resolver.Logger.Println(err.Error())
	}
	return response
}

//Syncs the changes from memory to the local cache file.
func (resolver *Resolver) Close() {
	resolver.Cache.Sync()
}