package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"repot/pkg/bot"
	"repot/pkg/edusms"
	"repot/pkg/model"
	"repot/pkg/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *Api) download(c *gin.Context) {
	dbot, err := bot.Instance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	str := c.Query("d")
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		dbot.SendSimple("failed to decode base64 string", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 string"})
		return
	}

	data := model.Data{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		dbot.SendSimple("Failed to unmarshal student data", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Srudent Data"})
		return
	}

	resp := gin.H{
		"id":           data.Student.ID,
		"admission_no": data.Student.AdmissionNo,
		"full_name":    data.Student.FullName,
		"url":          "https://llacademy.ng/student-view/" + strconv.Itoa(int(data.Student.ID)),
	}

	client := edusms.GetInstance()
	r, err := result.New(client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.Render(&data)
	if err != nil {
		resp["error"] = "failed to render result due to: " + err.Error()
		dbot.SendComplex(err.Error(), data.Student)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	byteFile, err := r.Generate()
	if err != nil {
		dbot.SendComplex(err.Error(), data.Student)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	fmt.Println(c.Request.Method)

	c.Header("Content-Disposition", "attachment; filename=result.pdf")
	c.Data(http.StatusOK, "application/pdf", byteFile)
}
