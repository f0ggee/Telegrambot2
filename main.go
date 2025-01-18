package main

import (
	"log"
	"os"
	_ "strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

// Глобальные переменные
var (
	bot       *tgbotapi.BotAPI
	userState = make(map[int64]string)
)

// Пример структуры для кнопок
type button struct {
	name string
	data string
}

// /Меню профиля
func profile() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "создать профиль", data: "create"},
		{name: "Посмотреть профиль", data: "check"},
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		buttons = append(buttons, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Главное меню (пример)
func startMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Подсчет калорий", data: "calorie"},
		{name: "Тренировка", data: "traine"},
		{name: "Профиль", data: "profile"},
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		buttons = append(buttons, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Меню «Тренировка»
func traineMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Тренировка: лёгкий уровень", data: "Light"},
		{name: "Тренировка: средний уровень", data: "Midle"},
		{name: "Тренировка: сложный уровень", data: "Hard"},
		{name: "Назад", data: "back"},
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		buttons = append(buttons, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func enlightenment() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Прокачка бицепса", data: "Bicepslight"},
		{name: "Прокачка передний части руку", data: "handle up"},
		{name: "Прокачка средний части руки", data: "handle midle"},
		{name: "Прокачка Задний части руки", data: "handle behind"},
		{name: "Прокачка Триципса", data: "updgrade triceps"},
		{name: "Прокачка бицепса", data: "Bicepslight"},
		{name: "назад", data: "back2"},
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		buttons = append(buttons, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)

}

func enlightenmentmidle() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Прокачка бицепса", data: "BicepslightM"},
		{name: "Прокачка передний части руку", data: "handle upM"},
		{name: "Прокачка средний части руки", data: "handle midleM"},
		{name: "Прокачка Задний части руки", data: "handle behindM"},
		{name: "Прокачка Триципса", data: "updgrade tricepsM"},
		{name: "Прокачка бицепса", data: "BicepslightM"},
		{name: "назад", data: "back3"},
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		buttons = append(buttons, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)

}

// Пример меню "help"

func main() {
	// Загружаем токен из .env (или откуда вам удобно)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not loaded (it's okay if you have token in another place)")
	}

	botToken := os.Getenv("TG_BOT_API")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	// Настраиваем получение обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}

	// Главный цикл обработки
	for update := range updates {
		// Обрабатываем колбэки от инлайн-кнопок
		if update.CallbackQuery != nil {
			callbacks(update)
			callbackcslight(update)
			callbackcsMidlet(update)
			continue
		}

		// Обрабатываем входящие сообщения
		if update.Message != nil {
			if update.Message.IsCommand() {
				commands(update)
			} else {

			}
		}
	}
}

func callbackcsMidlet(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	switch data {

	case "back3":
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "выборе тренировок ")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "handle upM":
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, "Это тренировки для дома *Основная часть:*\n   - Отжимания на трицепс: 3 подхода по 12 повтор")

		bot.Send(msg)

	}

}

func callbackcslight(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	switch data {
	case "back2":
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, "выборе тренировок ")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	}

}

// Функция обработки колбэков
func callbacks(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	switch data {

	case "traine":
		// Удаляем старое сообщение
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}

		// Выводим меню тренировок
		msg := tgbotapi.NewMessage(chatID, "Это список тренировок по уровням:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "back":
		// Удаляем старое сообщение (например, меню тренировок)
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}

		// Возвращаемся в главное меню
		msg := tgbotapi.NewMessage(chatID, "Вы в главном меню:")
		msg.ReplyMarkup = startMenu()
		sendMessage(msg)

	case "Light":
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали лёгкий уровень")
		msg.ReplyMarkup = enlightenment()
		sendMessage(msg)

	case "Midle":
		del := tgbotapi.NewDeleteMessage(chatID, messageID)
		if _, err := bot.Send(del); err != nil {
			log.Printf("Error deleting message: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, "вы выбрали средний уровень")
		msg.ReplyMarkup = enlightenmentmidle()
		sendMessage(msg)
		// Логика среднего уровня
	case "Hard":
		// Логика сложного уровня
		sendText(chatID, "Вы выбрали сложный уровень")
	}
}

func commands(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID, "Выберите действие")
		msg.ReplyMarkup = startMenu()
		sendMessage(msg)

	case "train":
		msg := tgbotapi.NewMessage(chatID, "Это список тренировок по уровню сложности:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "profile":
		msg := tgbotapi.NewMessage(chatID, "Выберите действия:")
		msg.ReplyMarkup = profile()
		sendMessage(msg)

	default:
		sendText(chatID, "Неизвестная команда: "+update.Message.Command())
	}
}

// Функция обработки обычного текста (не команды и не колбэка)

// Пример функции суммирования чисел из строки

// Обёртка для отправки простого текстового сообщения
func sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Обёртка для отправки любого Chattable-сообщения
func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Заглушка для сохранения данных (например, в БД)
