package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
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
			debugNotes := ``
			__dummy(debugNotes)
		} // end select
	}
}

func (bh *TGBotHandlers) execCommand(chatID int64, userName, cmd string) []ChatMessage {
	var result = make([]ChatMessage, 0, 2)
	var r ChatMessage
	switch {
	//resolve_1
	case strings.HasPrefix(cmd, "/resolve_"):
		t, err := bh.resolveTask(cutPrefix(cmd, "/resolve_"), chatID)
		if err != nil {
			r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf("Задача не на вас")}
			result = append(result, r)
		} else {
			r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf(`Задача "%s" выполнена`, t.task)}
			result = append(result, r)
			r = ChatMessage{chatID: t.authorID, msgText: fmt.Sprintf(`Задача "%s" выполнена @%s`, t.task, userName)}
			result = append(result, r)
		}

	case strings.HasPrefix(cmd, "/unassign_"):
		t, err := bh.unassignTask(cutPrefix(cmd, "/unassign_"), chatID)
		if err != nil {
			r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf("Задача не на вас")}
			result = append(result, r)
		} else {
			r = ChatMessage{chatID: t.authorID, msgText: fmt.Sprintf(`Задача "%s" осталась без исполнителя`, t.task)}
			result = append(result, r)
			r = ChatMessage{chatID: chatID, msgText: "Принято"}
			result = append(result, r)
		}

	case strings.HasPrefix(cmd, "/assign_"):
		// prevExecutorID, taskText, err := bh.assignTask(cutPrefix(cmd, "/assign_"), chatID, userName)
		prevExecutor, t, err := bh.assignTask(cutPrefix(cmd, "/assign_"), chatID, userName)
		panicOnError("assign failed", err)
		if prevExecutor == 0 {
			prevExecutor = t.authorID
		}
		r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf(`Задача "%s" назначена на вас`, t.task)}
		result = append(result, r)
		if t.authorID != chatID {
			r = ChatMessage{chatID: prevExecutor, msgText: fmt.Sprintf(`Задача "%s" назначена на @%s`, t.task, userName)}
			result = append(result, r)
		}

	case cmd == "/owner":
		tasks := bh.getTasksByAuthor(chatID)
		tasks = slices.DeleteFunc(tasks, func(x *UserTask) bool {
			show("filter: ", x)
			return x.executorID != chatID
		})
		if len(tasks) == 0 {
			r = ChatMessage{chatID: chatID, msgText: "Нет задач"}
		} else {
			var msgtxt = ""
			for _, t := range tasks {
				show("task: ", t)
				if len(msgtxt) > 0 {
					msgtxt += "\n\n"
				}
				s := "%s. %s by @%s\n/assign_%s"
				msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.taskID)
			}
			r = ChatMessage{chatID: chatID, msgText: msgtxt}
		}
		result = append(result, r)

	case cmd == "/my":
		tasks := bh.getTasksByExecutor(chatID)
		tasks = slices.DeleteFunc(tasks, func(x *UserTask) bool {
			show("filter: ", x)
			return x.authorID != chatID
		})
		if len(tasks) == 0 {
			r = ChatMessage{chatID: chatID, msgText: "Нет задач"}
		} else {
			var msgtxt = ""
			for _, t := range tasks {
				show("task: ", t)
				if len(msgtxt) > 0 {
					msgtxt += "\n\n"
				}
				s := "%s. %s by @%s\n/unassign_%s /resolve_%s"
				msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.taskID, t.taskID)
			}
			r = ChatMessage{chatID: chatID, msgText: msgtxt}
		}
		result = append(result, r)

	case cmd == "/tasks":
		// tasks := bh.getTasksByExecutorOrAuthor(chatID)
		// tasks = slices.DeleteFunc(tasks, func(x *UserTask) bool {
		// 	show("filter, task: ", x)
		// 	return x.resolved
		// })
		tasks := bh.getActiveTasks()
		if len(tasks) == 0 {
			r = ChatMessage{chatID: chatID, msgText: "Нет задач"}
		} else {
			var msgtxt = ""
			for _, t := range tasks {
				show("task: ", t)
				if len(msgtxt) > 0 {
					msgtxt += "\n\n"
				}
				if chatID != t.executorID && chatID == t.authorID {
					s := "%s. %s by @%s\nassignee: @%s"
					msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.executorName)
				} else if chatID == t.executorID && chatID != t.authorID {
					s := "%s. %s by @%s\nassignee: я\n/unassign_%s /resolve_%s"
					msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.taskID, t.taskID)
				} else if chatID == t.executorID && chatID == t.authorID && t.assigned {
					s := "%s. %s by @%s\nassignee: я\n/unassign_%s /resolve_%s"
					msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.taskID, t.taskID)
				} else {
					s := "%s. %s by @%s\n/assign_%s"
					msgtxt += fmt.Sprintf(s, t.taskID, t.task, t.authorName, t.taskID)
				}
			} // end tasks loop
			r = ChatMessage{chatID: chatID, msgText: msgtxt}
		}
		result = append(result, r)

	case strings.HasPrefix(cmd, "/new "):
		ut := bh.addTask(chatID, userName, cutPrefix(cmd, "/new "))
		r = ChatMessage{chatID: chatID, msgText: fmt.Sprintf(`Задача "%s" создана, id=%s`, ut.task, ut.taskID)}
		result = append(result, r)

	default:
		result = append(result, ChatMessage{chatID: chatID, msgText: "unknown command: " + cmd})
	}

	show("execCommand result: ", result)
	return result
}

func (bh *TGBotHandlers) getTaskByTaskID(taskID string) (idx int, task *UserTask, err error) {
	for idx, ut := range bh.data {
		if ut.taskID == taskID {
			return idx, ut, nil
		}
	}
	return 0, nil, fmt.Errorf("getTaskByTaskID, task %s not found", taskID)
}

func (bh *TGBotHandlers) getActiveTasks() []*UserTask {
	res := make([]*UserTask, 0, 16)
	for _, ut := range bh.data {
		if !ut.resolved {
			res = append(res, ut)
		}
	}
	return res
}

func (bh *TGBotHandlers) getTasksByExecutorOrAuthor(chatID int64) []*UserTask {
	res := make([]*UserTask, 0, 16)
	for _, ut := range bh.data {
		if ut.executorID == chatID || ut.authorID == chatID {
			res = append(res, ut)
		}
	}
	return res
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

func (bh *TGBotHandlers) getTasksByAuthor(chatID int64) []*UserTask {
	res := make([]*UserTask, 0, 16)
	for _, ut := range bh.data {
		if ut.authorID == chatID {
			res = append(res, ut)
		}
	}
	return res
}

func (bh *TGBotHandlers) resolveTask(taskID string, userID int64) (*UserTask, error) {
	idx, task, err := bh.getTaskByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("Can't resolve task that I can't find: %w", err)
	}
	log.Printf("task by id %s, found under index %d", taskID, idx)
	show("resolveTask: ", task)

	if task.executorID != userID {
		return task, fmt.Errorf("resolveTask, executor != user")
	} else {
		task.resolved = true
	}
	return task, nil
}

func (bh *TGBotHandlers) unassignTask(taskID string, userID int64) (*UserTask, error) {
	idx, task, err := bh.getTaskByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("Can't unassign task that I can't find: %w", err)
	}
	log.Printf("task by id %s, found under index %d", taskID, idx)
	show("unassignTask: ", task)

	if task.executorID != userID {
		return task, fmt.Errorf("unassignTask, executor != user")
	} else {
		task.executorID = 0
		task.executorName = ""
		task.assigned = false
	}
	return task, nil
}

func (bh *TGBotHandlers) assignTask(taskID string, newExecutorID int64, newExecutorName string) (int64, *UserTask, error) {
	idx, task, err := bh.getTaskByTaskID(taskID)
	if err != nil {
		return 0, nil, fmt.Errorf("Can't assign task that I can't find: %w", err)
	}
	log.Printf("task by id %s, found under index %d", taskID, idx)
	show("assignTask: ", task)

	prevExecutorID := task.executorID
	task.executorID = newExecutorID
	task.executorName = newExecutorName
	task.assigned = true

	return prevExecutorID, task, nil
}

func (bh *TGBotHandlers) addTask(authorID int64, authorName, task string) UserTask {
	ut := UserTask{
		taskID:       nextID_10(),
		task:         task,
		authorID:     authorID,
		authorName:   authorName,
		assigned:     false,
		executorID:   authorID,
		executorName: authorName,
		resolved:     false,
	}
	bh.data = append(bh.data, &ut)
	return ut
}

type UserTask struct {
	taskID       string
	task         string
	authorID     int64
	authorName   string
	assigned     bool
	executorID   int64
	executorName string
	resolved     bool
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
