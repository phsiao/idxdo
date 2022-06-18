package cmd

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
)

func colorPrettyJson(input interface{}) ([]byte, error) {
	buffer, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	json.Unmarshal(buffer, &obj)

	f := colorjson.NewFormatter()
	f.Indent = 2

	s, err := f.Marshal(obj)
	return s, err
}
