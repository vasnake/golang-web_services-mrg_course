package main

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"

	// "encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func init() {
	// upd global var for testing
	// we use patched version of gopkg.in/telegram-bot-api.v4 ( WebhookURL const -> var)
	WebhookURL = "http://127.0.0.1:8081"
	BotToken = "_golangcourse_test"
}

var (
	client = &http.Client{Timeout: time.Second}
)

// TDS is Telegram Dummy Server
type TDS struct {
	*sync.Mutex
	Answers map[int]string
}

func NewTDS() *TDS {
	return &TDS{
		Mutex:   &sync.Mutex{},
		Answers: make(map[int]string),
	}
}

func (srv *TDS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/getMe", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":{"id":` +
			strconv.Itoa(BotChatID) +
			`,"is_bot":true,"first_name":"game_test_bot","username":"game_test_bot"}}`))
	})
	mux.HandleFunc("/setWebhook", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":true,"description":"Webhook was set"}`))
	})
	mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		chatID, _ := strconv.Atoi(r.FormValue("chat_id"))
		text := r.FormValue("text")
		srv.Lock()
		srv.Answers[chatID] = text
		srv.Unlock()

		//fmt.Println("TDS sendMessage", chatID, text)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("unknown command %s", r.URL.Path))
	})

	handler := http.StripPrefix("/bot"+BotToken, mux)
	handler.ServeHTTP(w, r)
}

const (
	Ivanov     int = 256
	Petrov     int = 512
	Alexandrov int = 1024
	BotChatID      = 100500
)

var (
	users = map[int]*tgbotapi.User{
		Ivanov: &tgbotapi.User{
			ID:           Ivanov,
			FirstName:    "Ivan",
			LastName:     "Ivanov",
			UserName:     "ivanov",
			LanguageCode: "ru",
			IsBot:        false,
		},
		Petrov: &tgbotapi.User{
			ID:           Petrov,
			FirstName:    "Petr",
			LastName:     "Pertov",
			UserName:     "ppetrov",
			LanguageCode: "ru",
			IsBot:        false,
		},
		Alexandrov: &tgbotapi.User{
			ID:           Alexandrov,
			FirstName:    "Alex",
			LastName:     "Alexandrov",
			UserName:     "aalexandrov",
			LanguageCode: "ru",
			IsBot:        false,
		},
	}

	updID uint64
	msgID uint64
)

func SendMsgToBot(userID int, text string) error {
	// reqText := `{
	// 	"update_id":175894614,
	// 	"message":{
	// 		"message_id":29,
	// 		"from":{"id":133250764,"is_bot":false,"first_name":"Vasily Romanov","username":"rvasily","language_code":"ru"},
	// 		"chat":{"id":133250764,"first_name":"Vasily Romanov","username":"rvasily","type":"private"},
	// 		"date":1512168732,
	// 		"text":"THIS SEND FROM USER"
	// 	}
	// }`

	atomic.AddUint64(&updID, 1)
	myUpdID := atomic.LoadUint64(&updID)

	// better have it per user, but lazy now
	atomic.AddUint64(&msgID, 1)
	myMsgID := atomic.LoadUint64(&msgID)

	user, ok := users[userID]
	if !ok {
		return fmt.Errorf("no user %v", userID)
	}

	upd := &tgbotapi.Update{
		UpdateID: int(myUpdID),
		Message: &tgbotapi.Message{
			MessageID: int(myMsgID),
			From:      user,
			Chat: &tgbotapi.Chat{
				ID:        int64(user.ID),
				FirstName: user.FirstName,
				UserName:  user.UserName,
				Type:      "private",
			},
			Text: text,
			Date: int(time.Now().Unix()),
		},
	}
	reqData, _ := json.Marshal(upd)

	reqBody := bytes.NewBuffer(reqData)
	req, _ := http.NewRequest(http.MethodPost, WebhookURL, reqBody)
	_, err := client.Do(req)
	return err
}

type testCase struct {
	user    int
	command string
	answers map[int]string
}

func TestTasks(t *testing.T) {

	tds := NewTDS() // это мой эмулятор телеграм сервера
	ts := httptest.NewServer(tds)
	tgbotapi.APIEndpoint = ts.URL + "/bot%s/%s"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := startTaskBot(ctx, ":8081")
		if err != nil {
			panic(fmt.Sprintf("startTaskBot error: %s", err))
			// t.Fatalf("startTaskBot error: %s", err)
		}
	}()

	// give server time to start
	time.Sleep(10 * time.Millisecond)

	cases := []testCase{
		{
			// команда /tasks - выводит список всех активных задач
			Ivanov,
			"/tasks",
			map[int]string{
				Ivanov: "Нет задач",
			},
		}, // 0
		{
			// команда /new - создаёт новую задачу, всё что после /new - идёт в название задачи
			Ivanov,
			"/new написать бота",
			map[int]string{
				Ivanov: `Задача "написать бота" создана, id=1`,
			},
		}, // 1
		{
			Ivanov,
			"/tasks",
			map[int]string{
				Ivanov: `1. написать бота by @ivanov
/assign_1`,
			},
		}, // 2
		{
			// /assign_* - назначает задачу на себя
			Alexandrov,
			"/assign_1",
			map[int]string{
				Alexandrov: `Задача "написать бота" назначена на вас`,
				Ivanov:     `Задача "написать бота" назначена на @aalexandrov`,
			},
		}, // 3
		{
			// в случае если задача была назначена на кого-то - он получает уведомление об этом
			// в данном случае она была назначена на Alexandrov, поэтому ему отправляется уведомление
			Petrov,
			"/assign_1",
			map[int]string{
				Petrov:     `Задача "написать бота" назначена на вас`,
				Alexandrov: `Задача "написать бота" назначена на @ppetrov`,
			},
		}, // 4
		{
			// если задача назначена и на мне - показывается "на меня"
			Petrov,
			"/tasks",
			map[int]string{
				Petrov: `1. написать бота by @ivanov
assignee: я
/unassign_1 /resolve_1`,
			},
		}, // 5
		{
			// если задача назначена и не на мне - показывается логин исполнителя
			// при
			Ivanov,
			"/tasks",
			map[int]string{
				Ivanov: `1. написать бота by @ivanov
assignee: @ppetrov`,
			},
		}, // 6

		{
			// /unassign_ - снимает задачу с себя
			// нельзя снять задачу которая не на вас
			Alexandrov,
			"/unassign_1",
			map[int]string{
				Alexandrov: `Задача не на вас`,
			},
		}, // 7

		{
			// /unassign_ - снимает задачу с себя
			// автору отправляется уведомление что задача осталась без исполнителя
			Petrov,
			"/unassign_1",
			map[int]string{
				Petrov: `Принято`,
				Ivanov: `Задача "написать бота" осталась без исполнителя`,
			},
		}, // 8

		{
			// повтор
			// в случае если задача была назначена на кого-то - автор получает уведомление об этом
			Petrov,
			"/assign_1",
			map[int]string{
				Petrov: `Задача "написать бота" назначена на вас`,
				Ivanov: `Задача "написать бота" назначена на @ppetrov`,
			},
		}, // 9
		{
			// /resolve_* завершает задачу, удаляет её из хранилища
			// автору отправляется уведомление об этом
			Petrov,
			"/resolve_1",
			map[int]string{
				Petrov: `Задача "написать бота" выполнена`,
				Ivanov: `Задача "написать бота" выполнена @ppetrov`,
			},
		}, // 10

		{
			Petrov,
			"/tasks",
			map[int]string{
				Petrov: `Нет задач`,
			},
		}, // 11

		{
			// обратите внимание, id=2 - автоинкремент
			Petrov,
			"/new сделать ДЗ по курсу",
			map[int]string{
				Petrov: `Задача "сделать ДЗ по курсу" создана, id=2`,
			},
		}, // 12
		{
			// обратите внимание, id=3 - автоинкремент
			Ivanov,
			"/new прийти на хакатон",
			map[int]string{
				Ivanov: `Задача "прийти на хакатон" создана, id=3`,
			},
		}, // 13
		{
			Petrov,
			"/tasks",
			map[int]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
/assign_2

3. прийти на хакатон by @ivanov
/assign_3`,
			},
		}, // 14
		{
			// повтор
			// в случае если задача была назначена на кого-то - автор получает уведомление об этом
			// если он автор задачи - ему не приходит дополнительного уведомления о том что она назначена на кого-то
			Petrov,
			"/assign_2",
			map[int]string{
				Petrov: `Задача "сделать ДЗ по курсу" назначена на вас`,
			},
		}, // 15
		{
			Petrov,
			"/tasks",
			map[int]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
assignee: я
/unassign_2 /resolve_2

3. прийти на хакатон by @ivanov
/assign_3`,
			},
		}, // 16
		{
			// /my показывает задачи которые назначены на меня
			// при этому тут нет метки assegnee
			Petrov,
			"/my",
			map[int]string{
				Petrov: `2. сделать ДЗ по курсу by @ppetrov
/unassign_2 /resolve_2`,
			},
		}, // 17
		{
			// /owner - показывает задачи, которы я создал
			// при этому тут нет метки assegnee
			Ivanov,
			"/owner",
			map[int]string{
				Ivanov: `3. прийти на хакатон by @ivanov
/assign_3`,
			},
		}, // 18
	}

	for idx, item := range cases {

		tds.Lock()
		tds.Answers = make(map[int]string)
		tds.Unlock()

		caseName := fmt.Sprintf("[case%d, %d: %s]", idx, item.user, item.command)
		err := SendMsgToBot(item.user, item.command)
		if err != nil {
			t.Fatalf("%s SendMsgToBot error: %s", caseName, err)
		}
		// give TDS time to process request
		time.Sleep(10 * time.Millisecond)

		tds.Lock()
		result := reflect.DeepEqual(tds.Answers, item.answers)
		if !result {
			t.Fatalf("%s bad results:\n\tWant: %v\n\tHave: %v", caseName, item.answers, tds.Answers)
		}
		tds.Unlock()

	}

}
