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
	bot.Debug = false
	show("bot, authorized on account: ", bot.Self.UserName)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	panicOnError("SetWebhook failed", err)

	// info, err := bot.GetWebhookInfo()
	// panicOnError("bot.GetWebhookInfo failed", err)
	// if info.LastErrorDate != 0 {
	// 	show("Telegram callback failed: ", info.LastErrorMessage)
	// }

	// updates := bot.ListenForWebhook("/" + bot.Token)
	updates := bot.ListenForWebhook("/")
	// defer bot.StopReceivingUpdates()

	bh := NewTGBotHandlers(ctx, bot, updates)

	show("startTaskBot: program started ...")
	show("Starting server at: ", WebhookURL)
	go http.ListenAndServe(portStr, nil)

	return bh.ProcessMessages()
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

func (bh *TGBotHandlers) ProcessMessages() error {
	for {
		select {

		case <-bh.ctx.Done():
			show("ctx.Done")
			return nil

		case update := <-bh.updates:
			// show(fmt.Sprintf("update: %+v\n", update))
			if update.Message == nil { // ignore any non-Message Updates
				show("skip nil message: ", update)
				continue
			}

			id := update.Message.Chat.ID
			req := update.Message.Text
			show("bot got a new message, (chatID, text, user, messageID): ", id, req, update.Message.From.UserName, update.Message.MessageID)

			var msg = tgbotapi.NewMessage(id, "")
			msg.Text = bh.execCommand(id, req)
			// msg.ReplyToMessageID = update.Message.MessageID
			bh.bot.Send(msg)
			// _, err := bh.bot.Send(msg)
			// panicOnError("bot.Send failed", err)

			// [case#, user: command]
			debugNotes := `
2024-05-29T08:02:53.720Z: bot got a new message, (chatID, text, user, messageID): 256; "/tasks"; "ivanov"; 3; 
2024-05-29T08:02:53.720Z: execCommand result: "Нет задач";
	taskbot_test.go:390: [case2, 256: /tasks] bad results:
				Want: map[256:1. написать бота by @ivanov
`
			__dummy(debugNotes)
		} // end select
	}
}

func (bh *TGBotHandlers) execCommand(id int64, cmd string) string {
	// (id, text): 256; "/tasks";
	// "Нет задач"
	result := "unknown command: " + cmd
	switch {

	case cmd == "/tasks":
		result = "Нет задач"

	case strings.HasPrefix(cmd, "/new "):
		result = bh.addTask(id, cutPrefix(cmd, "/new "))

	default:
		result = "unknown command: " + cmd

	}

	show("execCommand result: ", result)
	return result
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
func __dummy(a any) error {
	if tgbotapi.APIEndpoint == a {
		return nil
	}
	return fmt.Errorf("xz")
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
