package utils

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
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

func UnmarshalJson(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	obj, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(obj, data)
	if err != nil {
		return err
	}

	return nil
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
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

func GetHomeDir() (string, error) {
	dir, ok := os.LookupEnv("SESSION_DIR")
	if ok {
		return filepath.Abs(dir)
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}

	return filepath.Abs(filepath.Join(dir, ".td"))
}
