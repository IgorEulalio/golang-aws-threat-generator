package events

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
	Action    string    `json:"Action"`
	Condition Condition `json:"Condition,omitempty"`
}

type Condition map[string]map[string]string

type Principal struct {
	Service string `json:"Service,omitempty"`
	AWS     string `json:"AWS,omitempty"`
}

func ParsePolicyFromEncodedString(encodedStr string) (Policy, error) {
	decodedStr, err := url.QueryUnescape(encodedStr)
	if err != nil {
		return Policy{}, fmt.Errorf("error decoding URL-encoded string: %w", err)
	}

	var policy Policy
	err = json.Unmarshal([]byte(decodedStr), &policy)
	if err != nil {
		return Policy{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	return policy, nil
}
