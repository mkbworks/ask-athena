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
	//Flag to enable or disable Trace logs
	traceLogs bool
}

// Queries the DNS server and fetches the 't' type record for 'name'.
func (resolver *Resolver) Resolve(name string, t RecordType) {
	MsgId := Id()
	resolver.response = NewMessage(MSG_RESOLVER_RESPONSE, MsgId)
	resolver.response.NewQuestion(name, t)
	if t == TYPE_A {
		_, err := resolver.resolveA(name)
		if err != nil {
			resolver.Log(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_AAAA {
		_, err := resolver.resolveAAAA(name)
		if err != nil {
			resolver.Log(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_TXT {
		_, err := resolver.resolveTXT(name)
		if err != nil {
			resolver.Log(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		fmt.Println(resolver.response.String())
	} else if t == TYPE_CNAME {
		_, err := resolver.resolveCNAME(name)
		if err != nil {
			resolver.Log(err.Error())
			resolver.response.Header.SetResponseCode(RC_SERVFAIL)
			fmt.Println(resolver.response.String())
			return
		}

		fmt.Println(resolver.response.String())
	} else {
		resolver.Log(ErrInvalidRecordType.Error())
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

//Adds the given resource records to resolver cache.
func (resolver *Resolver) addToCache(resources []Resource) {
	for _, RR := range resources {
		resolver.Cache.Add(RR.Name.Value, RR.TTL, RR.Class.String(), RR.Type.String(), RR.GetData())
	}
}

// Resolves the given domain name and returns its A resource records.
func (resolver *Resolver) resolveA(name string) ([]Resource, error) {
	cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_A)
	if ok {
		resolver.addToResolverResponse(name, cacheRecords)
		resolver.Log(fmt.Sprintf("A type records for %s have been served from the cache.", name))
		return cacheRecords, nil
	}

	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST, resolver.response.Header.Identifier)
	request.NewQuestion(name, TYPE_A)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
				resolver.addToCache(CNAME_RRs)
				return resolver.resolveA(CNAME_RRs[0].GetData())
			}

			A_RRs, _ := response.FindAnswerRecords(TYPE_A)
			resolver.addToResolverResponse(name, A_RRs)
			resolver.addToCache(A_RRs)
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
	cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_AAAA)
	if ok {
		resolver.addToResolverResponse(name, cacheRecords)
		resolver.Log(fmt.Sprintf("AAAA type records for %s have been served from the cache.", name))
		return cacheRecords, nil
	}

	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST, resolver.response.Header.Identifier)
	request.NewQuestion(name, TYPE_AAAA)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, exists := response.FindAnswerRecords(TYPE_CNAME)
			if exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
				resolver.addToCache(CNAME_RRs)
				return resolver.resolveAAAA(CNAME_RRs[0].GetData())
			}

			AAAA_RRs, _ := response.FindAnswerRecords(TYPE_AAAA)
			resolver.addToResolverResponse(name, AAAA_RRs)
			resolver.addToCache(AAAA_RRs)
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
	cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_TXT)
	if ok {
		resolver.addToResolverResponse(name, cacheRecords)
		resolver.Log(fmt.Sprintf("TXT type records for %s have been served from the cache.", name))
		return cacheRecords, nil
	}

	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST, resolver.response.Header.Identifier)
	request.NewQuestion(name, TYPE_TXT)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			TXT_RRs, _ := response.FindAnswerRecords(TYPE_TXT)
			resolver.addToResolverResponse(name, TXT_RRs)
			resolver.addToCache(TXT_RRs)
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

// Resolves the given domain name and returns the CNAME resource records.
func (resolver *Resolver) resolveCNAME(name string) ([]Resource, error) {
	cacheRecords, ok := resolver.Cache.Resolve(name, TYPE_CNAME)
	if ok {
		resolver.addToResolverResponse(name, cacheRecords)
		resolver.Log(fmt.Sprintf("CNAME type records for %s have been served from the cache.", name))
		return cacheRecords, nil
	}

	nameserver := resolver.getRootServer(TYPE_A)
	request := NewMessage(MSG_REQUEST, resolver.response.Header.Identifier)
	request.NewQuestion(name, TYPE_CNAME)
	for {
		response := resolver.getResponse(request, nameserver)
		if response.Header.AnCount > 0 {
			CNAME_RRs, Exists := response.FindAnswerRecords(TYPE_CNAME)
			if Exists {
				resolver.addToResolverResponse(name, CNAME_RRs)
				resolver.addToCache(CNAME_RRs)
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
	resolver.Log("**********************************************")
	resolver.Log(fmt.Sprintf("DNS Request being sent to server - %s.", ServerAddress))
	resolver.Log("**********************************************")
	resolver.Log(fmt.Sprintf("Request Contents are:\n%s", request.String()))
	resolver.Log("**********************************************")
	SendBuffer := request.Pack()
	udpConnect := UdpConnect{}
	err := udpConnect.ConnectTo(ServerAddress, DNS_PORT_NUMBER)
	if err != nil {
		resolver.Log(err.Error())
		return nil
	}
	var response *Message
	for validResponse := false; !validResponse; {
		err = udpConnect.Send(SendBuffer)
		if err != nil {
			resolver.Log(err.Error())
			return nil
		}
		receiveBuffer, err := udpConnect.Receive()
		if err != nil {
			resolver.Log(err.Error())
			return nil
		}
		response = NewMessage(MSG_RESPONSE, 0)
		response.Unpack(receiveBuffer)
		validResponse = response.IsResponse(request)
	}
	resolver.Log(fmt.Sprintf("Response received back:\n%s", response.String()))
	resolver.Log("**********************************************")
	err = udpConnect.Close()
	if err != nil {
		resolver.Log(err.Error())
	}
	return response
}

// Syncs the changes from memory to the local cache file.
func (resolver *Resolver) Close() {
	resolver.Cache.Sync()
}

//Logs information to the log file.
func (resolver *Resolver) Log(message string) {
	if resolver.traceLogs {
		resolver.Logger.Println(message)
	}
}
