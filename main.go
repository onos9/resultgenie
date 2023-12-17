package main

import (
	"fmt"
	"repot/pkg/api"
	"repot/pkg/workerpool"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	withErr := workerpool.WithErrorCallback(func(err error) {
		fmt.Println("Task error:", err)
	})

	pool := workerpool.New(100, withErr)
	defer pool.Release()

	api := api.New(router, pool)

	pool.Wait()
	api.Route()
	err := api.Run(":3000")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
