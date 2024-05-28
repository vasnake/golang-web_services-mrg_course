package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var WebhookURL = "http://127.0.0.1:8081"
var BotToken = "_golangcourse_test"

// err := startTaskBot(context.Background(), ":8081")
func startTaskBot(ctx context.Context, portStr string) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	panicOnError("NewBotAPI failed", err)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	panicOnError("SetWebhook failed", err)
	updates := bot.ListenForWebhook("/")

	bh := NewTGBotHandlers(ctx, bot, updates)

	show("startTaskBot: program started ...")
	show("Starting server at: ", WebhookURL)

	go bh.ProcessMessages()
	return http.ListenAndServe(portStr, nil)
}

type TGBotHandlers struct {
	ctx     context.Context
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewTGBotHandlers(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *TGBotHandlers {
	return &TGBotHandlers{
		ctx:     ctx,
		bot:     bot,
		updates: updates,
	}
}

func (bh *TGBotHandlers) ProcessMessages() {
	for {
		select {

		case <-bh.ctx.Done():
			show("ctx.Done")
			return

		case update := <-bh.updates:
			id := update.Message.Chat.ID
			cmd := update.Message.Text
			show("bot got a new message, (id, text): ", id, cmd)
			resp := bh.execCommand(id, cmd)
			show("bot sending response: ", resp)
			bh.bot.Send(tgbotapi.NewMessage(id, resp))

			// [case#, user: command]
			debugNotes := `
2024-05-28T14:39:01.582Z: bot got a new message, (id, text): 256; "/new написать бота";
2024-05-28T14:39:01.582Z: bot sending response: "Задача \"написать бота\" создана, id=1";
	taskbot_test.go:390: [case1, 256: /new написать бота] bad results:
				Want: map[256:Задача "написать бота" создана, id=1]
				Have: map[]
`
			__dummy(debugNotes)
		} // end select
	}
}

func (bh *TGBotHandlers) execCommand(id int64, cmd string) string {
	// (id, text): 256; "/tasks";
	// "Нет задач"
	switch {

	case cmd == "/tasks":
		return "Нет задач"

	case strings.HasPrefix(cmd, "/new "):
		return bh.addTask(id, cutPrefix(cmd, "/new "))

	default:
		show("unknown command: ", cmd)
		return "xz"
	}
}

func (bh *TGBotHandlers) addTask(id int64, task string) string {
	ut := UserTask{
		userID: id,
		taskID: nextID_10(),
		task:   task,
	}
	res := fmt.Sprintf(`Задача "%s" создана, id=%s`, ut.task, ut.taskID)
	// show("addTask result: ", res)
	return res
}

type UserTask struct {
	userID int64
	taskID string
	task   string
}

// это заглушка чтобы импорт сохранился
func __dummy(a any) {
	tgbotapi.APIEndpoint = "_dummy"
}

// --- useful little functions ---

var atomicCounter = new(atomic.Uint64)

func nextID_36() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 36)
}

func nextID_10() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 10)
}

func cutPrefix(s, prefix string) string {
	res, _ := strings.CutPrefix(s, prefix)
	return res
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
