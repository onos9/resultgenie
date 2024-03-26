package edusms

import (
	"encoding/json"
	"fmt"
	"io"
	"repot/pkg/model"
)

func (c *Client) GetStudentData(id string, data *model.Response) error {
	url := fmt.Sprintf("/marks-grade?id=%s", id)
	body, err := c.Get(url)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}
