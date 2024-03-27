package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"repot/pkg/bot"
	"repot/pkg/model"
	"repot/pkg/result"
	"repot/pkg/utils"
	"repot/pkg/workerpool"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DIR         = "generated"
	FILE_PREFIX = "repot-"
)

func (a *Api) cache(c *gin.Context) {
	dbot, err := bot.Instance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := model.Data{}
	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Student == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No exam data, please select Exam Type."})
		dbot.SendSimple("No Exam Data", "Please select Exam Type")
		return
	}

	id := strconv.Itoa(int(data.Student.ID))
	hash := md5.Sum([]byte(id))
	fileID := hex.EncodeToString(hash[:])
	filename := fmt.Sprintf("%s/%s%s.json", DIR, FILE_PREFIX, fileID)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		dbot.SendSimple("Failed to open file for writing: "+fileID, err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	byteFile, err := json.Marshal(data)
	if err != nil {
		dbot.SendComplex("Failed to write student data", err.Error(), data.Student)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = f.Write(byteFile)
	if err != nil {
		dbot.SendComplex("Failed to write PDF", err.Error(), data.Student)
		return
	}

	pool := workerpool.GetWorkerPool()
	pool.AddTask(func() (interface{}, error) {
		r, err := result.New()
		if err != nil {
			return nil, nil
		}
		byteFile, err = r.Render(&data)
		if err != nil {
			dbot.SendComplex("Failed to render PDF", err.Error(), data.Student)
			return nil, nil
		}

		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			dbot.SendComplex("Failed to open file for writing PDF", err.Error(), data.Student)
			return nil, err
		}
		defer f.Close()

		name := fmt.Sprintf("%s/%s%s.pdf", DIR, FILE_PREFIX, fileID)
		err = os.Rename(filename, name)
		if err != nil {
			dbot.SendComplex("Failed to rename file to .pdf", err.Error(), data.Student)
			return nil, err
		}
		_, err = f.Write(byteFile)
		if err != nil {
			dbot.SendComplex("Failed to write PDF", err.Error(), data.Student)
			return nil, err
		}

		return nil, nil
	})

	c.AbortWithStatus(http.StatusOK)
}

func (a *Api) download(c *gin.Context) {
	dbot, err := bot.Instance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	hash := md5.Sum([]byte(id))
	fileID := hex.EncodeToString(hash[:])
	filename := fmt.Sprintf("%s/%s%s.pdf", DIR, FILE_PREFIX, fileID)

	f, err := os.Open(filename)
	if err != nil {
		filename = fmt.Sprintf("%s/%s%s.json", DIR, FILE_PREFIX, fileID)
		f, err = os.Open(filename)
		if err != nil {
			dbot.SendSimple("Failed to open file for writing PDF, ID: "+id, err.Error())
			return
		}
	}
	defer f.Close()

	var byteFile []byte
	data := model.Data{}
	ext := filepath.Ext(f.Name())
	if ext == ".pdf" {
		byteFile, err = io.ReadAll(f)
		if err != nil {
			dbot.SendSimple("Failed to read PDF content, ID: "+id, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF content"})
			return
		}
	} else if ext == ".json" {
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&data)
		if err != nil {
			dbot.SendSimple("Failed to decode JSON data", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode JSON data"})
			return
		}

		resp := gin.H{
			"id":           data.Student.ID,
			"admission_no": data.Student.AdmissionNo,
			"full_name":    data.Student.FullName,
			"url":          "https://llacademy.ng/student-view/" + strconv.Itoa(int(data.Student.ID)),
			"code":         http.StatusOK,
			"error":        "",
		}

		r, err := result.New()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		byteFile, err = r.Render(&data)
		if err != nil {
			resp["error"] = err.Error()
			resp["code"] = http.StatusInternalServerError
			dbot.SendComplex("Failed to render result", err.Error(), data.Student)
			c.HTML(http.StatusInternalServerError, "error.html", resp)
			return
		}
	}

	time.AfterFunc(5*time.Minute, func() {
		err = utils.DeleteFile(filename)
		if err != nil {
			dbot.SendComplex("Failed to rename file to .pdf", err.Error(), data.Student)
			return
		}
		fmt.Println("File deleted: " + filename)
	})

	c.Header("Content-Disposition", "attachment; filename=result.pdf")
	c.Data(http.StatusOK, "application/pdf", byteFile)
}
