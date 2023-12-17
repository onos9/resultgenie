package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
)

func GetLocation(address string) (city, state string) {
	cityRegex := regexp.MustCompile(`,\s*([^,]+)$`)
	stateRegex := regexp.MustCompile(`,\s*([^,]+)\s*$`)

	cityMatches := cityRegex.FindStringSubmatch(address)
	stateMatches := stateRegex.FindStringSubmatch(address)

	if len(cityMatches) > 1 {
		city = strings.TrimSpace(cityMatches[1])
	}

	if len(stateMatches) > 1 {
		state = strings.TrimSpace(stateMatches[1])
	}

	return city, state
}

func UnmarshalJason(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	obj, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal(obj, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func EncodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func DecodeFromBase64(v interface{}, enc string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(enc))).Decode(v)
}

func StructToMap(obj interface{}) (newMap map[string]string, err error) {
	data, err := json.Marshal(obj) // Convert to a json strin
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	if err != nil {
		return
	}
	return
}
