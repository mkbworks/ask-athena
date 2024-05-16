package dns

import "errors"

//Identifies a protocol family or instance of a protocol.
type ClassType uint16

//Returns string representation of RR Class type.
func (ct ClassType) String() string {
	switch ct {
	case CLASS_IN:
		return "IN"
	case CLASS_CH:
		return "CH"
	default:
		return ""
	}
}

//Represents a record type in DNS
type RecordType uint16

//Returns string representation of the DNS RecordType. 
func (rt RecordType) String() string {
	switch rt {
	case TYPE_A:
		return "A"
	case TYPE_AAAA:
		return "AAAA"
	case TYPE_CNAME:
		return "CNAME"
	case TYPE_NS:
		return "NS"
	case TYPE_TXT:
		return "TXT"
	default:
		return ""
	}
}

//Maps a record type string  with its enumerated value.
type RecordTypes map[string]RecordType

//Fetches the key from RecordTypes instance, whose value matches 'recType'.
func (recTypes *RecordTypes) GetKey(recType RecordType) string {
	for key, value := range *recTypes {
		if value == recType {
			return key
		}
	}

	return ""
}

//Gets all the keys present in RecordTypes instance.
func (recTypes *RecordTypes) GetAllKeys() []string {
	keys := make([]string, 0)
	for key := range *recTypes {
		keys = append(keys, key)
	}
	return keys
}

//Gets the record type associated with the given string
func (recTypes RecordTypes) GetRecordType(key string) RecordType {
	recType, ok := recTypes[key]
	if !ok {
		panic(errors.New("record type not available"))
	}
	return recType
}

//Maps a class type string with its enumerated value.
type ClassTypes map[string]ClassType

//Gets the class type associated with the given string.
func (clsTypes ClassTypes) GetClassType(key string) ClassType {
	clsType, ok := clsTypes[key]
	if !ok {
		panic(errors.New("class type not available"))
	}
	return clsType
}