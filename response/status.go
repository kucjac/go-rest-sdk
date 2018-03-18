package response

import (
	"encoding/json"
	"strings"
)

// ResponseStatus is a basic status for API Response
// A developer can easily manage the response just by knowing the
// short status value.
// The status is a binary value - either the request where operated correctly or there exists error
type ResponseStatus int

const (
	Unknown ResponseStatus = iota
	StatusOk
	StatusError
)

// String - implements the Stringer interface
func (s ResponseStatus) String() string {
	switch s {
	case 1:
		return "ok"
	case 2:
		return "error"
	default:
		return "unknown"
	}
}

// MarshalJSON - implements json marshaller interface
func (s *ResponseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON - implements json Unmarshaler interface
func (s *ResponseStatus) UnmarshalJSON(b []byte) error {
	var status string
	if err := json.Unmarshal(b, &status); err != nil {
		return err
	}
	switch strings.ToLower(status) {
	case "ok":
		*s = 1
	case "error":
		*s = 2
	default:
		*s = 0
	}
	return nil
}
