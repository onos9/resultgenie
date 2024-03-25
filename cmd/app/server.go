package app

import (
	"fmt"
	"repot/pkg/edusms"
	"repot/pkg/workerpool"

	"go.uber.org/zap"
)

type App struct {
	pool workerpool.WorkerPool
	log  *zap.Logger
}

func New() *App {
	return &App{}
}

func (a *App) Run() {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()
	a.log = log

	_, err = edusms.New()
	if err != nil {
		panic(err)
	}

	withErr := workerpool.WithErrorCallback(func(err error) {
		fmt.Println("Task error:", err)
	})

	a.pool = workerpool.New(100, withErr)
	defer a.pool.Release()

	a.botServer()
	a.apiServer()

	a.pool.Wait()

}
