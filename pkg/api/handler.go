package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repot/pkg/auth"
	"repot/pkg/renderer"

	"github.com/gin-gonic/gin"
)

func (a *Api) process(c *gin.Context) {
	var data interface{}
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	a.Pool.AddTask(func() (interface{}, error) {
		d := data.(map[string]interface{})
		resultData, ok := d["data"]
		if !ok {
			return nil, fmt.Errorf("failed to get resultData data")
		}

		r, err := renderer.New(&auth.CLIENT)
		if err != nil {
			return nil, err
		}
		err = r.Render(resultData)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Task added to queue",
	})
}
