package dns

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//In-memory representation of a local resource stored in a BIND file.
type LocalResource struct {
	LastModified time.Time
	resource *Resource
}

//Returns the string representation of the local resource record.
func (lr *LocalResource) String() string {
	resourceString := lr.resource.CacheString()
	return fmt.Sprintf("%s%s%s\n", resourceString, WHITESPACE, lr.LastModified.Format(time.RFC3339))
}

//In-Memory representation of a BIND file.
type BindFile struct {
	//Resource records present in the BIND file.
	ResourceRecords []LocalResource
	//Local file path of the BIND file
	LocalFilePath string
}

//Initialize the attributes of BindFile instance.
func (bf *BindFile) Initialize(filePath string) error {
	bf.ResourceRecords = make([]LocalResource, 0)
	bf.LocalFilePath = filePath
	err := bf.Load()
	if err != nil {
		return err
	}

	return nil
}

//Creates a new local resource object and returns a pointer to the object.
func (bf *BindFile) NewLocalResource(name string, ttl uint32, class string, recType string, data string, LastModified string) *LocalResource {
	localResource := LocalResource{}
	localResource.resource = NewResourceRecord(name, ttl, class, recType, data)
	lastMod, err := time.Parse(time.RFC3339 ,LastModified)
	if err != nil {
		localResource.LastModified = time.Now().UTC()
	} else {
		localResource.LastModified = lastMod
	}
	return &localResource
}

//Creates a new local resource record and adds it to the BIND file if it has not already expired.
func (bf *BindFile) Add(name string, ttl uint32, class string, recType string, data string) {
	CurrentTime := time.Now().UTC()
	if ttl != 0 && !bf.HasRecordExpired(ttl, CurrentTime) {
		localResource := bf.NewLocalResource(name, ttl, class, recType, data, CurrentTime.Format(time.RFC3339))
		bf.ResourceRecords = append(bf.ResourceRecords, *localResource)
	}
}

//Load the RRs from the BIND file into memory.
func (bf *BindFile) Load() error {
	fileHandler, err := os.Open(bf.LocalFilePath)
	if err != nil {
		return err
	}
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	for {
		NewLine, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		NewLine = strings.TrimSuffix(NewLine, NEWLINE_SEPERATOR)
		NewLine = strings.TrimSpace(NewLine)
		if len(NewLine) != 0 {
			values := strings.Split(NewLine, WHITESPACE)
			if len(values) != 6 {
				return ErrParametersMissing
			} else {
				domainNameString := values[0]
				ttlString := values[1]
				ttlValue := uint32(parseUIntString(ttlString, 32))
				classString := values[2]
				typeString := values[3]
				dataString := values[4]
				lastModifiedString := values[5]
				newResource := bf.NewLocalResource(domainNameString, ttlValue, classString, typeString, dataString, lastModifiedString)
				bf.ResourceRecords = append(bf.ResourceRecords, *newResource)
			}
		}
	}

	return nil
}

//Persists the in-memory RR changes to the disk.
func (bf *BindFile) Sync() error {
	fileHandler, err := os.Create(bf.LocalFilePath)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileHandler)
	for _, rr := range bf.ResourceRecords {
		if !bf.HasRecordExpired(rr.resource.TTL, rr.LastModified) {
			writer.WriteString(rr.String())
		}
	}
	
	writer.Flush()
	fileHandler.Close()

	return nil
}

// Resolves the given domain name and record type using data available in the BIND file.
func (bf *BindFile) Resolve(name string, recType RecordType) ([]Resource, bool) {
	if recType == TYPE_A{
		return bf.resolveA(name)
	} else if recType == TYPE_AAAA {
		return bf.resolveAAAA(name)
	} else if recType == TYPE_CNAME {
		return bf.resolveCNAME(name)
	} else if recType == TYPE_TXT {
		return bf.resolveTXT(name)
	} else {
		return nil, false
	}
}

// Determines the A record for the given domain name from the BIND file.
func (bf *BindFile) resolveA(name string) ([]Resource, bool) {
	resources := make([]Resource, 0)
	CNAME_RRs, ok := bf.FindResources(name, TYPE_CNAME)
	if ok {
		A_RRs, ok := bf.resolveA(CNAME_RRs[0].GetData())
		if ok {
			resources = append(resources, CNAME_RRs...)
			resources = append(resources, A_RRs...)
		}
	}

	A_RRs, ok := bf.FindResources(name, TYPE_A)
	if ok {
		resources = append(resources, A_RRs...)
	}

	if len(resources) > 0 {
		return resources, true
	} else {
		return nil, false
	}
}

// Determines the AAAA record for the given domain name from the BIND file.
func (bf *BindFile) resolveAAAA(name string) ([]Resource, bool) {
	resources := make([]Resource, 0)
	CNAME_RRs, ok := bf.FindResources(name, TYPE_CNAME)
	if ok {
		AAAA_RRs, ok := bf.resolveAAAA(CNAME_RRs[0].GetData())
		if ok {
			resources = append(resources, CNAME_RRs...)
			resources = append(resources, AAAA_RRs...)
		}
	}

	AAAA_RRs, ok := bf.FindResources(name, TYPE_AAAA)
	if ok {
		resources = append(resources, AAAA_RRs...)
	}

	if len(resources) > 0 {
		return resources, true
	} else {
		return nil, false
	}
}

// Determines the CNAME record for the given domain name from the BIND file.
func (bf *BindFile) resolveCNAME(name string) ([]Resource, bool) {
	resources := make([]Resource, 0)
	CNAME_RRs, ok := bf.FindResources(name, TYPE_CNAME)
	if ok {
		resources = append(resources, CNAME_RRs...)
	}

	if len(resources) > 0 {
		return resources, true
	} else {
		return nil, false
	}
}

// Determines the TXT record for the given domain name from the BIND file.
func (bf *BindFile) resolveTXT(name string) ([]Resource, bool) {
	resources := make([]Resource, 0)
	TXT_RRs, ok := bf.FindResources(name, TYPE_TXT)
	if ok {
		resources = append(resources, TXT_RRs...)
	}

	if len(resources) > 0 {
		return resources, true
	} else {
		return nil, false
	}
}

//Returns all cached records matching the given domain name and record type.
func (bf *BindFile) FindResources(name string, recType RecordType) ([]Resource, bool) {
	resolvedValues := make([]Resource, 0)
	name = Canonicalize(name)
	if len(bf.ResourceRecords) > 0 {
		for _, lrr := range bf.ResourceRecords {
			if lrr.resource.Type == recType && strings.EqualFold(name, lrr.resource.Name.Value) && !bf.HasRecordExpired(lrr.resource.TTL, lrr.LastModified) {
				resolvedValues = append(resolvedValues, *lrr.resource)
			}
		}
	}

	if len(resolvedValues) == 0 {
		return resolvedValues, false
	}

	return resolvedValues, true
}

//Checks if the local resource is expired and returns true if it is and false if it has not expired.
func (bf *BindFile) HasRecordExpired(ttl uint32, LastModified time.Time) bool {
	TimeSinceLastMod := time.Now().UTC().Sub(LastModified)
	TimeInSeconds := TimeSinceLastMod.Seconds()
	if TimeInSeconds > float64(ttl) {
		return true
	} else {
		return false
	}
}