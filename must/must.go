package must

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func Url(r string) *url.URL {
	u, err := url.Parse(r)
	if err != nil {
		panic(err)
	}
	return u
}

func Serialize(data any) []byte {
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
