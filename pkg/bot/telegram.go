package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"repot/pkg/model"
	"strconv"
	"time"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

var client *gotgproto.Client
var username string

const (
	REGEX = `(\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}|[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}|\b\d{4}\b`
)

type TeleBot struct {
	*gotgproto.Client

	appID    int
	appHash  string
	botToken string
	username string
}

func Instance() (*TeleBot, error) {
	if client == nil {
		return nil, errors.New("Bot client not initialized")
	}

	return &TeleBot{
		Client:   client,
		username: username,
	}, nil
}

func New() (*TeleBot, error) {
	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		return nil, errors.Wrap(err, "APP_ID not set or invalid")
	}

	appHash, ok := os.LookupEnv("APP_HASH")
	if !ok {
		return nil, errors.New("no APP_HASH provided")
	}

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		return nil, errors.New("no BOT_TOKEN provided")
	}

	username, ok = os.LookupEnv("USERNAME")
	if !ok {
		return nil, errors.New("no BOT_TOKEN provided")
	}

	clientType := gotgproto.ClientType{
		BotToken: botToken,
	}

	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		fmt.Printf("Waiting for flood, dur: %d\n", wait.Duration)
	})

	ratelimiter := ratelimit.New(rate.Every(time.Millisecond*100), 30)
	opts := &gotgproto.ClientOpts{
		InMemory:    true,
		Session:     sessionMaker.SimpleSession(),
		Middlewares: []telegram.Middleware{waiter, ratelimiter},
		RunMiddleware: func(origRun func(ctx context.Context, f func(ctx context.Context) error) (err error), ctx context.Context, f func(ctx context.Context) (err error)) (err error) {
			return origRun(ctx, func(ctx context.Context) error {
				return waiter.Run(ctx, f)
			})
		},
	}

	client, err = gotgproto.NewClient(appID, appHash, clientType, opts)
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	return &TeleBot{
		Client:   client,
		appID:    appID,
		appHash:  appHash,
		botToken: botToken,
		username: username,
	}, nil
}

func (t *TeleBot) Run(ctx context.Context, log *zap.Logger) error {
	dispatcher := t.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("start", t.start))
	dispatcher.AddHandler(handlers.NewCommand("result", t.result))

	dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), t.buttonCallback))
	dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("btn1_"), t.btnFirst))
	dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("btn2_"), t.btnSecond))
	dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("btn3_"), t.btnThird))

	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, t.response), 1)

	fmt.Printf("client (@%s) has been started...\n", t.Self.Username)

	log.Info("Bot started!")
	err := t.Idle()
	if err != nil {
		log.Fatal("failed to start client:", zapcore.Field{
			Key:    "error",
			String: err.Error(),
		})
	}

	return nil
}

func (t *TeleBot) SendComplex(subject, msg string, data *model.Student) error {
	ctx := t.CreateContext()
	res := ctx.Sender.ResolveDomain(t.username)
	mkp := &tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "🔗 Student Profile",
						URL:  fmt.Sprintf("https://llacademy.ng/student-view/%d", int(data.ID)),
					},
				},
			},
		},
	}

	options := []styling.StyledTextOption{
		styling.Bold(fmt.Sprintf("%s\n\n", subject)),
		styling.Blockquote(fmt.Sprintf("Student Name: %s\nStudent ID: %d\nAdmission No: %d\n", data.FullName, int(data.ID), int(data.AdmissionNo))),
		styling.Code(msg),
	}

	res.Markup(mkp)
	res.StyledText(ctx, options...)
	return nil
}

func (t *TeleBot) SendSimple(subject, message string) error {
	ctx := t.CreateContext()

	res := ctx.Sender.ResolveDomain(t.username)
	options := []styling.StyledTextOption{
		styling.Bold(fmt.Sprintf("%s\n", subject)),
		styling.Code(message),
	}

	res.StyledText(ctx, options...)
	return nil
}
