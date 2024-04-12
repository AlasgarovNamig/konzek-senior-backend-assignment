package utils

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		Log("ERROR", fmt.Sprintf("Failed to marshal to JSON: %v", err))
		panic("json: " + err.Error())
	}
	return string(data)
}
