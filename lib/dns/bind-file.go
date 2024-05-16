package dns

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

//In-Memory representation of a BIND file.
type BindFile struct {
	//Resource records present in the BIND file.
	ResourceRecords []Resource
	//Local file path of the BIND file
	LocalFilePath string
	//Last Modified time for the file
	LastModifiedTime time.Time
}

//Initialize the attributes of BindFile instance.
func (bf *BindFile) Initialize(filePath string) {
	bf.ResourceRecords = make([]Resource, 0)
	bf.LocalFilePath = filePath
	bf.Load()
	bf.SetLastModified()
}

//Sets the last modified time for the BIND file.
func (bf *BindFile) SetLastModified() {
	fileInfo, err := os.Stat(bf.LocalFilePath)
	if err != nil {
		panic(err)
	}
	bf.LastModifiedTime = fileInfo.ModTime()
}

//Load the RRs from the BIND file into memory.
func (bf *BindFile) Load() {
	fileHandler, err := os.Open(bf.LocalFilePath)
	if err != nil {
		panic(err)
	}
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	for {
		NewLine, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		NewLine = strings.TrimSuffix(NewLine, NEWLINE_SEPERATOR)
		NewLine = strings.TrimSpace(NewLine)
		if len(NewLine) != 0 {
			values := strings.Split(NewLine, WHITESPACE)
			if len(values) != 5 {
				panic(errors.New("some parameters are missing in the resource record"))
			} else {
				domainNameString := values[0]
				ttlString := values[1]
				ttlValue := uint32(parseUIntString(ttlString, 32))
				classString := values[2]
				typeString := values[3]
				dataString := values[4]
				newResource := NewResourceRecord(domainNameString, ttlValue, classString, typeString, dataString)
				bf.ResourceRecords = append(bf.ResourceRecords, *newResource)
			}
		}
	}
}

//Persists the in-memory RR changes to the local copy file.
func (bf *BindFile) Sync() {
	fileHandler, err := os.Create(bf.LocalFilePath)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(fileHandler)
	for _, rr := range bf.ResourceRecords {
		writer.WriteString(rr.String())
	}
	
	writer.Flush()
	fileHandler.Close()
	bf.SetLastModified()
}