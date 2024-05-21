package action

import "encoding/json"

func Add(params json.RawMessage) (interface{}, error) {
	// Implement your add logic here
	// Example implementation:
	var p struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	return p.A + p.B, nil
}
