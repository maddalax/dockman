package kv

import (
	"encoding/json"
	"fmt"
)

func MustSerialize(data any) []byte {
	switch data.(type) {
	case []byte:
		return data.([]byte)
	case string:
		return []byte(data.(string))
	default:
		serialized, err := json.Marshal(data)
		if err != nil {
			return []byte(fmt.Sprintf("%v", data))
		}
		return serialized
	}
}
