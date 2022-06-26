package cmd

import (
	"encoding/hex"
	"encoding/json"
	"strings"

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

// IsEthereumAccount checks if the input can be a valid account
func IsEthereumAccount(input string) bool {
	if len(input) != 42 {
		return false
	}
	if !strings.HasPrefix(input, "0x") {
		return false
	}
	_, err := hex.DecodeString(input[2:])
	return err == nil
}
