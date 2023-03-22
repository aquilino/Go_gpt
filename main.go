package main

import (
	"context"
	"log"
	"strings"

	openai "github.com/0x9ef/openai-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	token      = "token_bot"
	gpt3APIKey = "token gpt3APIKey" 
)

func crear(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := update.Message.CommandArguments()
	ctx := context.Background()
	e := openai.New(gpt3APIKey)
	r, err := e.Completion(ctx, &openai.CompletionOptions{
		Prompt:    []string{text},
		MaxTokens: 2048,
		Model:     openai.ModelTextDavinci003,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	respuesta := r.Choices[0].Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, respuesta)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func editar(bot *tgbotapi.BotAPI, update tgbotapi.Update, input string) {
	text := update.Message.CommandArguments()
	ctx := context.Background()
	e := openai.New(gpt3APIKey)
	editResp, err := e.Edit(ctx, &openai.EditOptions{
		Model:       openai.ModelTextDavinci003,
		Input:       input, //update.Message.ReplyToMessage.Text,
		Instruction: text,
		N:           1,
		Temperature: 1,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	resp := editResp.Choices[0].Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func crear_imagen(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := update.Message.CommandArguments()
	e := openai.New(gpt3APIKey)
	ctx := context.Background()
	r, err := e.ImageCreate(ctx, &openai.ImageCreateOptions{
		Prompt:         strings.Join([]string{text}, ""),
		N:              2,
		Size:           "1024x1024",
		ResponseFormat: "url",
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	image_url := r.Data[0].Url
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, image_url)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Debug = true
	log.Printf("Cuenta autorizada a: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	
	if err != nil {
		log.Fatal(err)
		return
	}

	var lastUpdate tgbotapi.Update
	for update := range updates {
		if update.Message != nil {
			lastUpdate = update
		}
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if update.Message.Command() == "gpt" {
				crear(bot, lastUpdate)
			} else if update.Message.Command() == "gptc" {
				editar(bot, lastUpdate, update.Message.ReplyToMessage.Text)
			} else if update.Message.Command() == "gpti" {
				crear_imagen(bot, lastUpdate)
			}
		}
	}
}
