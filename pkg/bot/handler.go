package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"repot/pkg/edusms"
	"repot/pkg/model"
	"repot/pkg/workerpool"
	"strconv"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

func (t *TeleBot) start(ctx *ext.Context, update *ext.Update) error {
	user := update.EffectiveUser()
	text := fmt.Sprintf("Hello %s, I am %s, Wellcome to Lighthouse Leading Academy. Choose an option below to get started", user.FirstName, ctx.Self.FirstName)
	replyOpts := &ext.ReplyOpts{
		Markup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "🤔 Check Result",
							Data: []byte("cb_pressed"),
						},
						&tg.KeyboardButtonURL{
							Text: "🔗 Student Portal",
							URL:  "https://llacademy.ng/login",
						},
					},
				},
			},
		},
	}
	_, _ = ctx.Reply(update, text, replyOpts)

	// End dispatcher groups so that bot doesn't echo /start command usage
	return dispatcher.EndGroups
}

func (t *TeleBot) result(ctx *ext.Context, update *ext.Update) error {
	user := update.EffectiveUser()
	text := fmt.Sprintf("Hello %s, I am @%s, which of you child would like to get a report card?", user.FirstName, ctx.Self.FirstName)

	replyOpts := &ext.ReplyOpts{
		Markup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "Anna Bezallel",
							Data: []byte("cb_pressed"),
						},
					},
				},

				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "Grace Daniel",
							Data: []byte("cb_pressed"),
						},
					},
				},
			},
		},
	}
	_, _ = ctx.Reply(update, text, replyOpts)
	// End dispatcher groups so that bot doesn't echo /start command usage
	return dispatcher.EndGroups
}

func reply(ctx *ext.Context, username string) {

	res := ctx.Sender.ResolveDomain(username)

	options := []styling.StyledTextOption{
		styling.Bold("Ok result it is. Please enter your child's admission number\n\n"),
		styling.Code("You can find your child's admission number on their previouse report card. Please enter it"),
	}

	res.StyledText(ctx, options...)
}

func (t *TeleBot) btnFirst(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   false,
		QueryID: query.QueryID,
	})

	user := update.EffectiveUser()
	reply(ctx, user.Username)
	return nil
}

func (t *TeleBot) btnSecond(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   false,
		QueryID: query.QueryID,
	})

	user := update.EffectiveUser()
	reply(ctx, user.Username)
	return nil
}

func (t *TeleBot) btnThird(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   false,
		QueryID: query.QueryID,
	})

	user := update.EffectiveUser()
	reply(ctx, user.Username)
	return nil
}

func (t *TeleBot) buttonCallback(ctx *ext.Context, update *ext.Update) error {
	query := update.CallbackQuery
	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		Alert:   false,
		QueryID: query.QueryID,
	})

	user := update.EffectiveUser()
	res := ctx.Sender.ResolveDomain(user.Username)

	r := tg.ReplyMarkupClass(&tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "🔗 Student Portal",
						URL:  "https://llacademy.ng/login",
					},
				},
			},
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "🔗 Student Portal",
						URL:  "https://llacademy.ng/login",
					},
				},
			},
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "🔗 Student Portal",
						URL:  "https://llacademy.ng/login",
					},
				},
			},
		},
	})

	res.Markup(r)

	options := []styling.StyledTextOption{
		styling.Bold("Ok result it is. Please enter your child's admission number\n\n"),
		styling.Code("You can find your child's admission number on their previouse report card. Please enter it"),
	}

	res.StyledText(ctx, options...)
	return nil
}

func (t *TeleBot) response(ctx *ext.Context, update *ext.Update) error {
	msg := update.EffectiveMessage

	val, ok := match(msg.Text, `\b\d{4}\b`)
	if !ok {
		return fmt.Errorf("invalid admission number")
	}

	admino, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error converting to integer:", err)
		return err
	}

	c := edusms.GetInstance()
	body, err := c.Get("/exam-list/0/0", nil)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	data := model.ExamData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	pool := workerpool.GetWorkerPool()
	for _, e := range data.ExamTypes {
		pool.AddTask(func() (interface{}, error) {
			c := edusms.GetInstance()
			url := fmt.Sprintf("/marks-grade?id=%d&exam_type=%d", admino, int(e.ID))
			body, err := c.Get(url, nil)
			if err != nil {
				return nil, err
			}

			b, err := io.ReadAll(body)
			if err != nil {
				return nil, err
			}

			data := []model.Data{}
			err = json.Unmarshal(b, &data)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})
	}

	return nil
}
