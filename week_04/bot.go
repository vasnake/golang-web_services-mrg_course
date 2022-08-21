package main

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/telegram-bot-api.v4" // you need to install that package
	"io/ioutil"
	"net/http"
)

const (
	BotToken   = "310805560:AAENzjDSJPKABY9Hw1GZOdKBxxrhOHkfo_k" // see BotFather in telegram
	WebhookURL = "https://ea731f5c.ngrok.io"                     // ok for testing
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	// unmarshal xml from rss using tags
	Items []Item `xml:"channel>item"`
}
type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	// get data from ext service
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil) // push to background
	fmt.Println("start listen :8080")

	// получаем все обновления из канала updates
	for update := range updates {
		if url, ok := rss[update.Message.Text]; ok { // we have that rss?
			rss, err := getNews(url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"sorry, error happend",
				))
			}

			for _, item := range rss.Items {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					item.URL+"\n"+item.Title,
				))
			}
		} else { // no, we have not that rss
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				`there is only Habr feed availible`,
			))
		}

	}

}
