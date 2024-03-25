package app

import (
	"context"
	"repot/pkg/bot"
)

func (a *App) botServer() {

	tg, err := bot.New()
	if err != nil {
		panic("[Error] failed to create Telegram client due to: " + err.Error())
	}

	a.pool.AddTask(func() (interface{}, error) {
		err = tg.Run(context.Background(), a.log)
		if err != nil {
			panic("[Error] failed to start Telegram Bot due to: " + err.Error())
		}

		return nil, nil
	})
}
