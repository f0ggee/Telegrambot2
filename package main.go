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
	bot       *tgbotapi.BotAPIk
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
		{name: "назад", data: "back_to_levels"},
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
		{name: "назад", data: "back_to_levels"},
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
func enlightenmentHard() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Прокачка бицепса", data: "BicepsHard"},
		{name: "Прокачка передний части руку", data: "handle uphard"},
		{name: "Прокачка средний части руки", data: "handle midlehard"},
		{name: "Прокачка Задний части руки", data: "handle behindhard"},
		{name: "Прокачка Триципса", data: "updgrade tricepshard"},
		{name: "Прокачка бицепса", data: "Bicepslighthard"},
		{name: "назад", data: "back_to_levels"},
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

		//file := tgbotapi.NewDocumentUpload(chatID, tgbotapi.FilePath("C:\\Users\\USER\\Desktop\\Traine\\Midlehandleup.txt"))

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

	// Обработка команды "start"
	switch update.Message.Command() {
	case "start":
		// Приветственное сообщение
		welcomeText := `Привет! Я ваш помощник-бот. Вот что я умею:
- 📸 Отправлять картинки.
- 🏋️‍♂️ Показывать список тренировок.
- 🧑‍💻 Профиль с полезной информацией.

Используйте меню ниже или введите команды вручную. Например:
/start - Начать работу
/train - Список тренировок
/profile - Профиль`

		// Приветственное сообщение
		msg := tgbotapi.NewMessage(chatID, welcomeText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send welcome message: %v", err)
			return
		}

		// Клавиатура для главного экрана
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("📸 Отправить картинку"),
				tgbotapi.NewKeyboardButton("🏋️‍♂️ Список тренировок"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🧑‍💻 Профиль"),
			),
		)
		keyboard.ResizeKeyboard = true // Подгонка кнопок

		// Отправка клавиатуры
		menuMsg := tgbotapi.NewMessage(chatID, welcomeText)
		menuMsg.ReplyMarkup = keyboard
		if _, err := bot.Send(menuMsg); err != nil {
			log.Printf("Failed to send menu: %v", err)
		}

	case "train":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень сложности для тренировок:")
		msg.ReplyMarkup = traineMenu()
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}

	case "profile":
		msg := tgbotapi.NewMessage(chatID, "Выберите действия:")
		msg.ReplyMarkup = profile()
		sendMessage(msg)

	default:
		sendText(chatID, "Неизвестная команда: "+update.Message.Command())
	}
}

// Обработка callback
func handleCallback(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	data := update.CallbackQuery.Data

	switch data {
	case "Light":
		msg := tgbotapi.NewMessage(chatID, "Выберите упражнение для лёгкого уровня:")
		msg.ReplyMarkup = enlightenment()
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending Light menu: %v", err)
		}

	case "Midle":
		msg := tgbotapi.NewMessage(chatID, "Выберите упражнение для среднего уровня:")
		msg.ReplyMarkup = enlightenmentMidle()
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending Midle menu: %v", err)
		}

	case "Hard":
		msg := tgbotapi.NewMessage(chatID, "Выберите упражнение для сложного уровня:")
		msg.ReplyMarkup = enlightenmentHard()
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send Hard menu: %v", err)
		}

	default:
		sendText(chatID, "Извините, я не понял ваш выбор.")
	}

	// Ответ на CallbackQuery
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

// Отправка текста
func sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Отправка сообщения
func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Заглушка для сохранения данных (например, в БД)
