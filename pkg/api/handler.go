package api

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"repot/pkg/bot"
	"repot/pkg/edusms"
	"repot/pkg/model"
	"repot/pkg/result"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *Api) cache(c *gin.Context) {
	data := map[string]string{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, ok := data["student_data"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_data not found"})
		return
	}

	dbot, err := bot.Instance()
	hash := md5.Sum([]byte(str))
	file_id := hex.EncodeToString(hash[:])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	f, err := os.Create(file_id)
	if err != nil {
		dbot.SendSimple("Failed to unmarshal student data", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		dbot.SendSimple("Failed to unmarshal student data", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id": file_id,
	})

}

func (a *Api) download(c *gin.Context) {
	dbot, err := bot.Instance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Query("id")
	examId := c.Query("exam_id")
	client := edusms.GetInstance()

	d := model.Response{}
	err = client.GetStudentData(id, examId, &d)
	if err != nil {
		dbot.SendSimple("Failed to get student data", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	data := d.Data
	resp := gin.H{
		"id":           data.Student.ID,
		"admission_no": data.Student.AdmissionNo,
		"full_name":    data.Student.FullName,
		"url":          "https://llacademy.ng/student-view/" + strconv.Itoa(int(data.Student.ID)),
	}

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

	c.Header("Content-Disposition", "attachment; filename=result.pdf")
	c.Data(http.StatusOK, "application/pdf", byteFile)
}
