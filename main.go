package main

import (
	con "MSB/config"
	mod "MSB/modules"
	"fmt"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var devs Devs

func main() {
	devs = make(Devs)
	devs["Денис"] = "@Jestric"
	devs["Александр"] = "@Kobalt_KSPA"
	devs["Глеб"] = "@Kavantek"

	config := mod.CreateParamServer("./config/config.json")
	mod.CheckParam(config, true)

	bot, err := tgbotapi.NewBotAPI(config.BotParams.Token)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	chatId := config.ChatID

	go Remind(config, bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatId = update.Message.Chat.ID
		fmt.Printf("ChatId: %v, Message: %v\n", chatId, update.Message.Text)
		go Router(update, bot)
	}
}

func Router(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	iventMatch, _ := regexp.MatchString("^Событие", update.Message.Text)

	fmt.Printf("ivent Trigered - %v\n", iventMatch)

	if iventMatch {
		go IventTimer(update, bot)
	}
	switch update.Message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я Бот-напоминалка. Буду отправлять тебе сообщения о твоих задачах!")
		bot.Send(msg)
	}
}

func Remind(config con.Config, bot *tgbotapi.BotAPI) {
	for {
		nowTime := time.Now().Hour()

		fmt.Printf("Time now: %v\n", nowTime)

		if nowTime < 8 || nowTime > 15 {
			time.Sleep(30 * time.Minute)
			continue
		}

		fmt.Printf("Trigered time now: %v\n", nowTime)

		today := fmt.Sprintf("%02d-%02d", time.Now().Day(), time.Now().Month())
		weekday := time.Now().Weekday()

		for i := range config.ReportPlanDays {
			if today == config.ReportPlanDays[i] {
				msg := tgbotapi.NewMessage(config.ChatID, `@Jestric @Kavantek @Kobalt_KSPA Конец месяца! Не забудьте заполнить ЛКО и сдать Отчёт "План Управление"!`)
				bot.Send(msg)
			}
		}

		for i := range config.NewMonthDays {
			if today == config.NewMonthDays[i] {
				msg := tgbotapi.NewMessage(config.ChatID, `@Jestric @Kavantek @Kobalt_KSPA Начало месяца! Не забудьте распечатать ЛКО и сформировать Отчёт программиста!`)
				bot.Send(msg)
			}
		}

		if fmt.Sprintf("%v", weekday) != "Sunday" && fmt.Sprintf("%v", weekday) != "Saturday" {
			nowDuration := time.Duration(time.Now().Hour())*time.Hour + time.Duration(time.Now().Minute())*time.Minute
			trigerTime := time.Hour*15 + time.Minute*50
			sleepTime := trigerTime - nowDuration
			fmt.Printf("Sleep time: %v", sleepTime)
			time.Sleep(sleepTime)
			msg := tgbotapi.NewMessage(config.ChatID, `@Jestric @Kavantek @Kobalt_KSPA 10 минут до сдачи ежедневного отчёта!`)
			bot.Send(msg)
			time.Sleep(time.Hour * 12)
		}
	}
}

func IventTimer(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	respMsg := `%v, напоминаю: через 15 минут у вас событие "%v"`
	splitMsg := strings.Split(update.Message.Text, "; ")

	members := strings.Split(splitMsg[1], ", ")

	for i := range members {
		members[i] = devs[members[i]]
	}

	curr, _ := time.Parse("02-01-2006 15:04:05", getNow())

	// Форматируем время события
	timeT, _ := time.Parse("02-01-2006 15:04:05", splitMsg[3])

	// Вычитаем из времени события время текущее
	delta := timeT.Sub(curr)

	delta = delta - (time.Minute * 15)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(`Создал событие "%v"! Напомню о нём за 15 минут до начала.`, splitMsg[2]))

	bot.Send(msg)

	time.Sleep(delta)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(respMsg, strings.Join(members, ", "), splitMsg[2]))

	bot.Send(msg)
}

// Возвращает текущее время в формате 11-05-2023 08:47:15
func getNow() string {
	return fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d", time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
}

type Devs map[string]string
