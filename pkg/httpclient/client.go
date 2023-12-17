package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAX_RETRY = 3
const RETRY_INTERVAL = 5 * time.Second

type Client struct {
	*http.Client
	Token   string
	host    string
	headers map[string]string
}

type Chunk struct {
	StreamID   string        `json:"stream_id"`
	ChunkIndex int           `json:"chunk_index"`
	Data       []interface{} `json:"data"`
}

func init() {
	os.Setenv("SERVER_HOST", "https://llacademy.ng")
}

func New() Client {
	headers := map[string]string{"Content-Type": "application/json"}
	c := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
	}
	// godotenv.Load(".env")
	return Client{
		Client:  c,
		host:    os.Getenv("SERVER_HOST"),
		headers: headers,
	}
}

func (c *Client) Send(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", c.Token)
	response, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to API endpoint. %+v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", response.Status)
	}

	return response, nil
}

func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *Client) Get(url string, payload *bytes.Buffer) (io.ReadCloser, error) {
	if payload == nil {
		payload = bytes.NewBuffer([]byte{})
	}
	url = fmt.Sprintf("%s%s", c.host, url)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s", err)
	}
	response, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c *Client) Post(url string, payload *bytes.Buffer) (io.ReadCloser, error) {
	if payload == nil {
		payload = bytes.NewBuffer([]byte{})
	}

	url = fmt.Sprintf("%s%s", c.host, url)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("Client: could not create request: %s", err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	response, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c *Client) GetWithStream(url string, payload *bytes.Buffer, onData func(chunk []interface{})) error {
	if payload == nil {
		payload = bytes.NewBuffer([]byte{})
	}

	receivedChunks := make(map[string][]int)
	for attempt := 1; attempt <= MAX_RETRY; attempt++ {
		url = fmt.Sprintf("%s%s", c.host, url)
		req, err := http.NewRequest("GET", url, payload)
		if err != nil {
			fmt.Println("Error creating request:", err)
			time.Sleep(RETRY_INTERVAL)
			continue
		}

		req.Header.Set("Authorization", c.Token)

		response, err := c.Do(req)
		if err != nil {
			fmt.Printf("error sending request to API endpoint. %+v", err)
			time.Sleep(RETRY_INTERVAL)
			continue
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			fmt.Printf("error response from server: %s", response.Status)
			time.Sleep(RETRY_INTERVAL)
			continue
		}

		totalChunksHeader := response.Header.Get("X-Total-Chunks")
		totalChunks, err := strconv.Atoi(totalChunksHeader)
		fmt.Println("totalChunks:", totalChunks)
		if err != nil {
			fmt.Println("Error converting totalChunks header:", err)
			return err
		}

		decoder := json.NewDecoder(response.Body)
		for {
			var chunk Chunk
			if err = decoder.Decode(&chunk); err == io.EOF {
				return nil
			} else if err != nil {
				fmt.Printf("error decoding JSON: %v", err)
				break
			}

			if _, exists := receivedChunks[chunk.StreamID]; !exists {
				receivedChunks[chunk.StreamID] = make([]int, 0)
			}

			if contains(receivedChunks[chunk.StreamID], chunk.ChunkIndex) {
				fmt.Printf("Received duplicate chunk: StreamID=%s, ChunkIndex=%d\n", chunk.StreamID, chunk.ChunkIndex)
				continue
			}

			onData(chunk.Data)

			receivedChunks[chunk.StreamID] = append(receivedChunks[chunk.StreamID], chunk.ChunkIndex)
		}

		if allChunksReceived(receivedChunks, totalChunks) {
			fmt.Println("All expected data received successfully.")
			break
		}
	}
	return nil
}

func (c *Client) CreateForm(form map[string]string, buf *bytes.Buffer) (*multipart.Writer, error) {
	mp := multipart.NewWriter(buf)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			name := val[strings.LastIndex(val, "/")+1:]
			part, err := mp.CreateFormFile(key, name)
			if err != nil {
				return nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp, nil
}

func allChunksReceived(receivedChunks map[string][]int, totalChunks int) bool {
	for _, indices := range receivedChunks {
		if len(indices) != totalChunks {
			return false // Some chunks are missing
		}
	}

	return true // All chunks are received
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
