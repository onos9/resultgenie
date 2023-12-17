package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"repot/pkg/httpclient"
)

var CLIENT httpclient.Client

func init() {
	postBody, _ := json.Marshal(map[string]string{
		"email":    "onosbrown.saved@gmail.com",
		"password": "#1414bruno#",
	})

	CLIENT = httpclient.New()
	body, err := CLIENT.Post("/api/auth/login", bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("LOGIN: could not login: %s\n", err)
	}

	buffer, err := io.ReadAll(body)
	if err != nil {
		log.Fatalf("LOGIN: could not read body: %s\n", err)
	}

	var auth map[string]interface{}
	err = json.Unmarshal(buffer, &auth)
	if err != nil {
		log.Fatalf("LOGIN: could not get token: %s\n", err)
	}

	data := auth["data"].(map[string]interface{})
	CLIENT.Token = data["token"].(string)
}
