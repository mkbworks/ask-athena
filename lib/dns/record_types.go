package dns

type ClassType uint16
type RecordType uint16

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