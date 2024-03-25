package edusms

import (
	"fmt"
	"os"
	"repot/pkg/httpclient"

	"github.com/joho/godotenv"
)

var client *httpclient.HTTPClient

type Client struct {
	*httpclient.HTTPClient
}

func New() (Client, error) {
	client = httpclient.New()
	c := Client{
		HTTPClient: client,
	}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	email, ok := os.LookupEnv("EMAIL")
	if !ok {
		panic("EMAIL not set")
	}
	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		panic("PASSWORD not set")
	}

	token, err := c.login(email, password)
	if err != nil {
		panic("could not login: " + err.Error())
	}

	c.SetHeader("Authorization", token)
	return c, nil
}

func GetInstance() *Client {
	return &Client{
		HTTPClient: client,
	}
}
