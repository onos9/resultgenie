package edusms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Auth struct {
	Email    string
	Password string
	Token    string
}

func (c *Client) login(email, password string) (string, error) {

	postBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	c.SetHeader("Content-Type", "application/json")
	body, err := c.Post("/auth/login", bytes.NewBuffer(postBody))
	if err != nil {
		return "", fmt.Errorf("[LOGIN]: could not login: %s", err)
	}

	buffer, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("LOGIN: could not read body: %s", err)
	}

	var auth map[string]interface{}
	err = json.Unmarshal(buffer, &auth)
	if err != nil {
		return "", fmt.Errorf("LOGIN: could not unmarshal body: %s", err)
	}

	data := auth["data"].(map[string]interface{})
	token := data["token"].(string)

	return token, nil
}
