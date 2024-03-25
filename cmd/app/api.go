package app

import (
	"repot/pkg/api"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *App) apiServer() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	a.pool.AddTask(func() (interface{}, error) {
		api := api.New(router, a.pool)
		api.Route()
		err := api.Run(":3000")
		if err != nil {
			panic("[Error] failed to start Gin server due to: " + err.Error())
		}

		return nil, nil
	})
}
