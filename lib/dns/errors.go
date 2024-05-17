package dns

import (
	"errors"
)

var ErrInvalidRecordType error = errors.New("record type requested is not allowed by the resolver")
var ErrNameServerFetch error = errors.New("unable to fetch Name Server details")
var ErrAuthNameServerFetch error = errors.New("unable to fetch Authoritative Name Server details")
var ErrMessageTooLong error = errors.New("udp message size is too long")
var ErrNotAbsolutePath error = errors.New("file path must be an absolute")
var ErrParametersMissing error = errors.New("parameters are missing")
var ErrBitCount error = errors.New("bit count for the given number is larger than the required bit count")