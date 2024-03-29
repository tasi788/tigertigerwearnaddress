package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func fetchContent(text string) []string {
	contentURL := "https://raw.githubusercontent.com/hexUniverse/postergirl/master/" + text
	resp, err := http.Get(contentURL)
	if err != nil {
		log.Fatalf("Fetch Error")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	content := strings.Split(string(body), ",")
	return content
}

func main() {
	token := os.Getenv("BOTTOKEN")
	admins := []int64{525239263, 184805205, 54465600}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Fetch Content
	dict := fetchContent("content.txt")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "update":
				for _, id := range admins {
					if id == update.Message.From.ID {
						dict = fetchContent("content.txt")
						msg.Text = "已更新。"
					}
				}

			default:
				continue
			}
			if msg.Text != "" {
				bot.Send(msg)
			}

		}

		if update.InlineQuery != nil {
			if checker(update.InlineQuery.From.ID) == true && update.InlineQuery.Query == "sk" {
				dict := fetchContent("rich.txt")
				article := tgbotapi.NewInlineQueryResultArticle(
					update.InlineQuery.ID,
					"台幣戰士",
					dict[0])
				article.ThumbURL = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/120/apple/325/eggplant_1f346.png"
				inlineConf := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					CacheTime:     1,
					IsPersonal:    false,
					Results:       []interface{}{article},
				}
				bot.Send(inlineConf)
				continue
			}

			rand.Seed(time.Now().Unix())
			article := tgbotapi.NewInlineQueryResultArticle(
				update.InlineQuery.ID,
				"虎虎？",
				dict[rand.Intn(len(dict))])

			article.ThumbURL = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/120/apple/198/dress_1f457.png"

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				CacheTime:     1,
				IsPersonal:    false,
				Results:       []interface{}{article},
			}
			bot.Send(inlineConf)
		}
	}
}

func checker(id int64) bool {
	fuckingrich := []int64{525239263, 1013801267}
	for _, richid := range fuckingrich {
		if id == richid {
			return true
		}
	}
	return false
}
