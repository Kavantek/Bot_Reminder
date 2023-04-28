package main

import (
	mod "MSB/modules"
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	config := mod.CreateParamServer("./config/config.json")
	mod.CheckParam(config, true)

	bot, err := tgbotapi.NewBotAPI(config.BotParams.Token)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		nowTime := time.Now().Hour()
		for nowTime < 8 {
			time.Sleep(30 * time.Minute)
		}
		today := fmt.Sprintf("%02d-%02d", time.Now().Day(), time.Now().Month())
		for i := range config.ReportPlanDays {
			if today == config.ReportPlanDays[i] {
				msg := tgbotapi.NewMessage(config.ChatID.Gleb, `Конец месяца! Не забудь заполнить ЛКО и сдать Отчёт "План Управление"!`)
				bot.Send(msg)
				msg = tgbotapi.NewMessage(config.ChatID.Alexander, `Конец месяца! Не забудь заполнить ЛКО и сдать Отчёт "План Управление"!`)
				bot.Send(msg)
				msg = tgbotapi.NewMessage(config.ChatID.Denis, `Конец месяца! Не забудь заполнить ЛКО и сдать Отчёт "План Управление"!`)
				bot.Send(msg)
			}
		}
		for i := range config.NewMonthDays {
			if today == config.NewMonthDays[i] {
				msg := tgbotapi.NewMessage(config.ChatID.Gleb, `Начало месяца! Не забудь распечатать ЛКО и сформировать Отчёт программиста!`)
				bot.Send(msg)
				msg = tgbotapi.NewMessage(config.ChatID.Alexander, `Начало месяца! Не забудь распечатать ЛКО и сформировать Отчёт программиста!`)
				bot.Send(msg)
				msg = tgbotapi.NewMessage(config.ChatID.Denis, `Начало месяца! Не забудь распечатать ЛКО и сформировать Отчёт программиста!`)
				bot.Send(msg)
			}
		}
	}()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		fmt.Printf("ChatId: %v, Message: %v\n", chatID, update.Message.Text)
		go Router(update, bot)
	}
}

func Router(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.Message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я Бот-напоминалка. Буду отправлять тебе сообщения о твоих задачах!")
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, я просто напоминалка и больше ничего не умею. Но если это изменится, я тебе об этом скажу!")
		bot.Send(msg)
	}
}
