package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var botToken = flag.String("token", "", "Bot Token")
var chatId = flag.Int64("chat", -1002200729105, "Bot Token")

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

func send(title string, id string) error {
	body, err := json.Marshal(map[string]any{
		//"parse_mode": ModeHTML,
		"text":    "тут має бути фотка",
		"chat_id": *chatId,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": title + ":" + id,
						"url":  "https://t.me/kormy6hkaBot/shop?startapp=" + id,
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

	if err := send("Черешня - Жовта", "1"); err != nil {
		panic(err)
	}
	//send("Персик")
}
