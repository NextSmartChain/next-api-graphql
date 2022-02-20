// Package types implements different core types of the API.
package types

import "encoding/json"

// ValidatorInfo holds extended validator information.
type ValidatorInfo struct {
	// Name represents the name of the validator
	Name *string `json:"name"`

	// LogoUrl represents validator logo URL
	LogoUrl *string `json:"logoUrl"`

	// Website represents a link to validator website
	Website *string `json:"website"`

	// Contact represents a link to contact to the validator
	Contact *string `json:"contact"`
}

// UnmarshalValidatorInfo parses the JSON-encoded validator information data.
func UnmarshalValidatorInfo(data []byte) (*ValidatorInfo, error) {
	var sfci ValidatorInfo
	err := json.Unmarshal(data, &sfci)
	return &sfci, err
}

// Marshal returns the JSON encoding of validator information.
func (sfci *ValidatorInfo) Marshal() ([]byte, error) {
	return json.Marshal(sfci)
}
