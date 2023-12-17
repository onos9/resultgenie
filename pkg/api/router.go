package api

import (
	"net/http"
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
	a.POST("/", a.process)

	a.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
