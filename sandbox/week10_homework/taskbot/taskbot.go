package main

import (
	"context"
	"fmt"
	"log"
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
	data    []*UserTask
}

func NewTGBotHandlers(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *TGBotHandlers {
	return &TGBotHandlers{
		ctx:     ctx,
		bot:     bot,
		updates: updates,
		data:    make([]*UserTask, 0, 16),
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

			chatID := update.Message.Chat.ID
			cmd := update.Message.Text
			userName := update.Message.From.UserName
			show("bot got a new message, (chatID, text, user, messageID): ", chatID, cmd, userName, update.Message.MessageID)

			var results = bh.execCommand(chatID, userName, cmd)
			for _, r := range results {
				var msg = tgbotapi.NewMessage(r.chatID, "")
				msg.Text = r.msgText
				// msg.ReplyToMessageID = update.Message.MessageID
				bh.bot.Send(msg)
				// _, err := bh.bot.Send(msg)
				// panicOnError("bot.Send failed", err)
			}

			// [case#, user: command]
			debugNotes := `
2024-05-29T09:49:14.539Z: bot got a new message, (chatID, text, user, messageID): 512; "/tasks"; "ppetrov"; 6;
2024-05-29T09:49:14.539Z: execCommand result: []main.ChatMessage{main.ChatMessage{chatID:512, msgText:"1. написать бота by @ivanov\n/assign_1"}};
	taskbot_test.go:390: [case5, 512: /tasks] bad results:

				Want: map[512:1. написать бота by @ivanov
		assignee: я
		/unassign_1 /resolve_1]

				Have: map[512:1. написать бота by @ivanov
		/assign_1]
`
			__dummy(debugNotes)
		} // end select
	}
}

func (bh *TGBotHandlers) execCommand(chatID int64, userName, cmd string) []ChatMessage {
	var result = make([]ChatMessage, 0, 2)
	var r ChatMessage
	switch {

	case cmd == "/tasks":
		tasks := bh.getTasksByExecutor(chatID)
		if len(tasks) == 0 {
			r = ChatMessage{chatID: chatID, msgText: "Нет задач"}
		} else {
			t := tasks[0]
			r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf("%s. написать бота by @%s\n/assign_%s", t.taskID, t.authorName, t.taskID)}
		}
		result = append(result, r)

	case strings.HasPrefix(cmd, "/new "):
		ut := bh.addTask(chatID, userName, cutPrefix(cmd, "/new "))
		r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf(`Задача "%s" создана, id=%s`, ut.task, ut.taskID)}
		result = append(result, r)

	case strings.HasPrefix(cmd, "/assign_"):
		prevExecutorID, taskText, err := bh.assignTask(cutPrefix(cmd, "/assign_"), chatID)
		panicOnError("assign failed", err)
		r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf(`Задача "%s" назначена на вас`, taskText)}
		result = append(result, r)
		r = ChatMessage{chatID: prevExecutorID, msgText: fmt.Sprintf(`Задача "%s" назначена на @%s`, taskText, userName)}
		result = append(result, r)

	default:
		result = append(result, ChatMessage{chatID: chatID, msgText: "unknown command: " + cmd})
	}

	show("execCommand result: ", result)
	return result
}

func (bh *TGBotHandlers) assignTask(taskID string, targetChatID int64) (prevExecutorID int64, taskText string, err error) {
	// find taks #1
	// assign new owher to task
	idx, task, err := bh.getTaskByTaskID(taskID)
	if err != nil {
		return 0, "", fmt.Errorf("Can't assign task that I can't find: %w", err)
	}
	log.Printf("task by id %s, found under index %d", taskID, idx)

	prevExecutor := task.executorID
	task.executorID = targetChatID

	return prevExecutor, task.task, nil
}

func (bh *TGBotHandlers) getTaskByTaskID(taskID string) (idx int, task *UserTask, err error) {
	for idx, ut := range bh.data {
		if ut.taskID == taskID {
			return idx, ut, nil
		}
	}
	return 0, nil, fmt.Errorf("getTaskByTaskID, task %s not found", taskID)
}

func (bh *TGBotHandlers) getTasksByExecutor(chatID int64) []*UserTask {
	res := make([]*UserTask, 0, 16)
	for _, ut := range bh.data {
		if ut.executorID == chatID {
			res = append(res, ut)
		}
	}
	return res
}

func (bh *TGBotHandlers) addTask(authorID int64, authorName, task string) UserTask {
	ut := UserTask{
		executorID: authorID,
		authorID:   authorID,
		taskID:     nextID_10(),
		authorName: authorName,
		task:       task,
	}
	bh.data = append(bh.data, &ut)
	return ut
}

type UserTask struct {
	executorID int64
	authorID   int64
	taskID     string
	authorName string
	task       string
}

type ChatMessage struct {
	chatID  int64
	msgText string
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
