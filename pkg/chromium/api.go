package chromium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"repot/pkg/httpclient"
	"strconv"
)

type Options struct {
	PaperWidth        float64 `json:"paperWidth"`
	PaperHeight       float64 `json:"paperHeight"`
	MarginTop         float64 `json:"marginTop"`
	MarginBottom      float64 `json:"marginBottom"`
	MarginLeft        float64 `json:"marginLeft"`
	MarginRight       float64 `json:"marginRight"`
	PreferCssPageSize bool    `json:"preferCssPageSize"`
	PrintBackground   bool    `json:"printBackground"`
	OmitBackground    bool    `json:"omitBackground"`
	Landscape         bool    `json:"landscape"`
	Scale             float64 `json:"scale"`
	NativePageRanges  string  `json:"nativePageRanges"`
}

// Chromium struct represents the Chromium module in Gotenberg.
type Chromium struct {
	options *Options
	*httpclient.HTTPClient
}

// NewChromium initializes a new Chromium instance.
func New(client *httpclient.HTTPClient) *Chromium {

	return &Chromium{
		options: &Options{
			MarginTop:    0.2,
			MarginBottom: 0.2,
			MarginLeft:   0.2,
			MarginRight:  0.2,
		},
		HTTPClient: client,
	}
}

// Request generates an HTTP request with the configured options.
func (c *Chromium) SendHTML(path string) ([]byte, error) {
	form := map[string]string{
		"files":        path,
		"marginTop":    "0.2",
		"marginBottom": "0.2",
		"marginLeft":   "0.2",
		"marginRight":  "0.2",
	}

	body := new(bytes.Buffer)
	w, err := c.CreateForm(form, body)
	if err != nil {
		return nil, err
	}

	url := os.Getenv("CHROMIUM_URL")
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	respons, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer respons.Body.Close()

	if respons.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", respons.StatusCode)
	}

	data, err := io.ReadAll(respons.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o Options) MarshalJSON() ([]byte, error) {
	type Alias Options
	return json.Marshal(&struct {
		*Alias
		PaperWidth        string `json:"paperWidth"`
		PaperHeight       string `json:"paperHeight"`
		MarginTop         string `json:"marginTop"`
		MarginBottom      string `json:"marginBottom"`
		MarginLeft        string `json:"marginLeft"`
		MarginRight       string `json:"marginRight"`
		PreferCssPageSize string `json:"preferCssPageSize"`
		PrintBackground   string `json:"printBackground"`
		OmitBackground    string `json:"omitBackground"`
		Landscape         string `json:"landscape"`
		Scale             string `json:"scale"`
		NativePageRanges  string `json:"nativePageRanges"`
	}{
		Alias:             (*Alias)(&o),
		PaperWidth:        fmt.Sprintf("%f", o.PaperWidth),
		PaperHeight:       fmt.Sprintf("%f", o.PaperHeight),
		MarginTop:         fmt.Sprintf("%f", o.MarginTop),
		MarginBottom:      fmt.Sprintf("%f", o.MarginBottom),
		MarginLeft:        fmt.Sprintf("%f", o.MarginLeft),
		MarginRight:       fmt.Sprintf("%f", o.MarginRight),
		PreferCssPageSize: strconv.FormatBool(o.PreferCssPageSize),
		PrintBackground:   strconv.FormatBool(o.PrintBackground),
		OmitBackground:    strconv.FormatBool(o.OmitBackground),
		Landscape:         strconv.FormatBool(o.Landscape),
		Scale:             fmt.Sprintf("%f", o.Scale),
		NativePageRanges:  o.NativePageRanges,
	})
}
