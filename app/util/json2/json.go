package json2

import "encoding/json"

func Deserialize[T any](data []byte) (*T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return &result, err
}

func SerializeOrEmpty[T any](data T) []byte {
	serialized, err := json.Marshal(data)
	if err != nil {
		return []byte{}
	}
	return serialized
}

func Serialize[T any](data T) ([]byte, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}
