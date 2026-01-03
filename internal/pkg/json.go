package pkg

import "encoding/json"

func ToJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}

func JsonToByte(data string) []byte {
	if data == "" {
		return []byte("{}")
	}
	return []byte(data)
}
