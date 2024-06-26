package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var botToken = flag.String("token", "", "Bot Token")
var chatId = flag.Int64("chat", -1002200729105, "Bot Token")
var botName = flag.String("botName", "kormy6hkaBot", "Bot Token")

type ParseMode string

const (
	ModeHTML     ParseMode = "HTML"
	ModeMarkdown ParseMode = "Markdown"
)

type ErrorResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func send(image, title string, id string) error {
	raw, _ := json.Marshal(map[string]any{
		"id":    id,
		"title": title,
	})

	body, err := json.Marshal(map[string]any{
		//"parse_mode": ModeHTML,
		"text":    image,
		"chat_id": *chatId,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": title,
						"url":  fmt.Sprintf("https://t.me/%s/shop?startapp=%s", *botName, base64.URLEncoding.EncodeToString(raw)),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("json.Marshal :%w", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *botToken),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("http.Post :%w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	answerByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll :%w", err)
	}

	var telegramError ErrorResponse
	if err = json.Unmarshal(answerByte, &telegramError); err != nil {
		return fmt.Errorf("json.Unmarshal :%w", err)
	}

	if !telegramError.OK {
		return fmt.Errorf("telegram response code:%d, mgs:%s",
			telegramError.ErrorCode,
			telegramError.Description,
		)
	}

	return nil
}

func main() {
	flag.Parse()

	if err := send("https://eimg.pravda.com/images/doc/2/1/2106f80-cherry-1566779.jpg", "Черешня - Жовта", "1"); err != nil {
		panic(err)
	}
	if err := send("https://beautyseason.ru/upload/iblock/371/325gbduwwkdko0hbu0bqcjqxa4jfvcks.jpg", "Персик", "2"); err != nil {
		panic(err)
	}
}
