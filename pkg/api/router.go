package api

import (
	"fmt"
	"net/http"
	"repot/pkg/bot"
	"repot/pkg/model"
	"repot/pkg/workerpool"

	"github.com/gin-gonic/gin"
)

type Api struct {
	*gin.Engine
	Pool workerpool.WorkerPool
}

func New(r *gin.Engine, pool workerpool.WorkerPool) Api {
	return Api{
		Engine: r,
		Pool:   pool,
	}
}

func (a *Api) Route() {

	a.GET("/download", a.download)
	a.POST("/:id", a.cache)

	a.GET("/bot", func(c *gin.Context) {
		dbot, err := bot.Instance()
		if err != nil {
			fmt.Printf("failed to create bot: %s\n", err.Error())
		}

		data := model.Student{
			ID:          1,
			FullName:    "John Doe",
			AdmissionNo: 1234,
		}

		dbot.SendComplex("api route called", data)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	a.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

}
