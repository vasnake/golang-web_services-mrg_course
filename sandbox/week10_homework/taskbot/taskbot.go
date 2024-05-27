package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var WebhookURL = "http://127.0.0.1:8081"
var BotToken = "_golangcourse_test"

// err := startTaskBot(context.Background(), ":8081")
func startTaskBot(ctx context.Context, portStr string) error {
	// сюда писать код
	/*
		в этом месте вы стартуете бота,
		стартуете хттп сервер который будет обслуживать этого бота
		инициализируете ваше приложение
		и потом будете обрабатывать входящие сообщения
	*/
	bot, err := tgbotapi.NewBotAPI(BotToken)
	panicOnError("NewBotAPI failed", err)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	panicOnError("SetWebhook failed", err)

	updates := bot.ListenForWebhook("/")

	show("startTaskBot: program started ...")
	show("Starting server at: ", WebhookURL)
	go http.ListenAndServe(portStr, nil)

	for update := range updates {
		id := update.Message.Chat.ID
		cmd := update.Message.Text
		show("bot got a new message, (id, text): ", id, cmd)
		bot.Send(tgbotapi.NewMessage(id, "Нет задач"))
	}

	// [case#, user: command]
	debugNotes := `
bot got a new message, (id, text): 256; "/new написать бота";
taskbot_test.go:390: [case1, 256: /new написать бота] bad results:
			Want: map[256:Задача "написать бота" создана, id=1]
			Have: map[256:Нет задач]
	`
	__dummy(debugNotes)

	return nil
}

type TGBotHandlers struct {
	ctx context.Context
}

func NewTGBotHandlers(ctx context.Context) *TGBotHandlers {
	return &TGBotHandlers{
		ctx: ctx,
	}
}

// ServeHTTP implements http.Handler.
func (serv *TGBotHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("ServeHTTP: ", r.URL, r.Method, r)
	http.Error(w, "oops 1", http.StatusNotImplemented)
}

var _ http.Handler = NewTGBotHandlers(nil) // type check

// это заглушка чтобы импорт сохранился
func __dummy(a any) {
	tgbotapi.APIEndpoint = "_dummy"
}

// --- useful little functions ---

var atomicCounter = new(atomic.Uint64)

func nextID() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 36)
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func strRef(in string) *string {
	return &in
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
