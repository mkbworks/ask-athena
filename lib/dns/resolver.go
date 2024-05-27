package dns

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
)

// Structure to represent a DNS Resolver.
type Resolver struct {
	//References the BIND file containing the DNS root server details.
	RootServers BindFile
	//References the BIND file containing all the cached resource records.
	Cache BindFile
	//Logger to be used to generate logs.
	Logger *log.Logger
	//References the DNS response being formed during domain name resolution.
	response *Message
}

// Queries the DNS server and fetches the 't' type record for 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) {
	resolver.response = NewMessage(MSG_RESOLVER_RESPONSE)
	resolver.response.NewQuestion(name, t)
	if t == TYPE_A {
		cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_A)
		if ok {
			resolver.Logger.Printf("IPv4 address for the domain name %s has been served from local Cache file.\n", name)
			resolver.response.AddAnswers(cacheRecords)
			fmt.Println(resolver.response.String())
			return
		}
		_, err := resolver.resolveA(name)
		if err != nil {
			resolver.Logger.Println(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		for _, A_RR := range resolver.response.Answers {
			resolver.Cache.Add(A_RR.Name.Value, A_RR.TTL, A_RR.Class.String(), A_RR.Type.String(), A_RR.GetData())
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_AAAA {
		cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_AAAA)
		if ok {
			resolver.Logger.Printf("IPv6 address for the domain name %s has been served from local Cache file.\n", name)
			resolver.response.AddAnswers(cacheRecords)
			fmt.Println(resolver.response.String())
			return
		}
		_, err := resolver.resolveAAAA(name)
		if err != nil {
			resolver.Logger.Println(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		for _, AAAA_RR := range resolver.response.Answers {
			resolver.Cache.Add(AAAA_RR.Name.Value, AAAA_RR.TTL, AAAA_RR.Class.String(), AAAA_RR.Type.String(), AAAA_RR.GetData())
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_TXT {
		cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_TXT)
		if ok {
			resolver.Logger.Printf("TXT record for the domain name %s has been served from local Cache file.\n", name)
			resolver.response.AddAnswers(cacheRecords)
			fmt.Println(resolver.response.String())
			return
		}
		_, err := resolver.resolveTXT(name)
		if err != nil {
			resolver.Logger.Println(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		for _, TXT_RR := range resolver.response.Answers {
			resolver.Cache.Add(TXT_RR.Name.Value, TXT_RR.TTL, TXT_RR.Class.String(), TXT_RR.Type.String(), TXT_RR.GetData())
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_CNAME {
		cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_CNAME)
		if ok {
			resolver.Logger.Printf("CNAME record for the domain name %s has been served from the local cache file.\n", name)
			resolver.response.AddAnswers(cacheRecords)
			fmt.Println(resolver.response.String())
			return
		}

		_, err := resolver.resolveCNAME(name)
		if err != nil {
			resolver.Logger.Println(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		for _, CNAME_RR := range resolver.response.Answers {
			resolver.Cache.Add(CNAME_RR.Name.Value, CNAME_RR.TTL, CNAME_RR.Class.String(), CNAME_RR.Type.String(), CNAME_RR.GetData())
		}

		fmt.Println(resolver.response.String())
	} else {
		resolver.Logger.Println(ErrInvalidRecordType.Error())
		resolver.response.Header.SetResponseCode(RC_NOTIMP)
		fmt.Println(resolver.response.String())
	}
}

// Adds the given resource records to resolver response message if the domain name given is present in the Question.
func (resolver *Resolver) addToResolverResponse(name string, resources []Resource) {
	if resolver.response.HasQuestion(name) {
		resolver.response.AddAnswers(resources)
	} else if resolver.response.HasAnswer(name) {
		resolver.response.AddAnswers(resources)
	}
}

// Resolves the given domain name and returns its A resource records.
func (resolver *Resolver) resolveA(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, TYPE_A)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
				return resolver.resolveA(CNAME_RRs[0].GetData())
			}

			A_RRs, _ := response.FindAnswerRecords(TYPE_A)
			resolver.addToResolverResponse(name, A_RRs)
			return  A_RRs, nil
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

// Resolves the given domain name and returns its AAAA resource records.
func (resolver *Resolver) resolveAAAA(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, TYPE_AAAA)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
				return resolver.resolveAAAA(CNAME_RRs[0].GetData())
			}

			AAAA_RRs, _ := response.FindAnswerRecords(TYPE_AAAA)
			resolver.addToResolverResponse(name, AAAA_RRs)
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

// Resolves the given domain name and returns its TXT resource records.
func (resolver *Resolver) resolveTXT(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, TYPE_TXT)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			TXT_RRs, _ := response.FindAnswerRecords(TYPE_TXT)
			resolver.addToResolverResponse(name, TXT_RRs)
			return TXT_RRs, nil
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

func (resolver *Resolver) resolveCNAME(name string) ([]Resource, error) {
	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST)
	request.NewQuestion(name, TYPE_CNAME)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, Exists := response.FindAnswerRecords(TYPE_CNAME)
			if Exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
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

// Returns true if the record type provided is accepted by the resolver, else returns false.
func (resolver *Resolver) IsAllowed(recordType string) bool {
	_, exists := AllowedRRTypes[recordType]
	return exists
}

// Returns the record type object for the given type string.
func (resolver *Resolver) GetRecordType(recordType string) RecordType {
	return AllowedRRTypes.GetRecordType(recordType)
}

// Returns the IP address of a root DNS server chosen at random.
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

// Sends the request to the target DNS server and receives a response over the same connection.
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

// Syncs the changes from memory to the local cache file.
func (resolver *Resolver) Close() {
	resolver.Cache.Sync()
}
