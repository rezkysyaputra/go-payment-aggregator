package pkg

import "encoding/json"

func ToJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func JsonToByte(data string) []byte {
	if data == "" {
		return []byte("{}")
	}
	return []byte(data)
}
