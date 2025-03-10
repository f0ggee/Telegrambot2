package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)


const (
	profileStepHeight = 1
	profileStepWeight = 2
	profileStepAge    = 3
	profileStepGender = 4
	profileStepDone   = 5
)

const (
	diaryStepAdd = 201
)

var userStep = make(map[int64]int)

type user_profile struct {
	Height int
	Weight int
	Age    int
	Gender string
}

var userProfiles = make(map[int64]*user_profile)

// Структура для описания «кнопок»
type button struct {
	name string
	data string
}

var bot *tgbotapi.BotAPI


func calculateCalories(gender string, weight float64, height float64, age int) float64 {
	if gender == "мужской" {
		return 88.36 + (13.4 * weight) + (4.8 * height) - (5.7 * float64(age))
	}
	return 447.6 + (9.2 * weight) + (3.1 * height) - (4.3 * float64(age))
}


type DiaryEntry struct {
	Date time.Time // Когда была сделана запись
	Text string    // Текст записи
}

var userDiary = make(map[int64][]DiaryEntry)

func diaryMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Добавить запись", data: "diary_add"},
		{name: "Показать записи", data: "diary_show"},
		{name: "Назад", data: "back"},
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		btn := tgbotapi.NewInlineKeyboardButtonData(st.name, st.data)
		row := tgbotapi.NewInlineKeyboardRow(btn)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func main() {
	// Загружаем токен из .env (или прописываем вручную)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not loaded (если токен прописан вручную, то всё ок)")
	}

	botToken := os.Getenv("TG_BOT_API")
	if botToken == "" {
		log.Fatal("Нет токена TG_BOT_API в .env")
	}

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram Bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Не удалось получить канал обновлений: %v", err)
	}

	log.Println("Бот запущен...")

	// Цикл обработки входящих сообщений/колбэков
	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(update)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommands(update)
			} else {
				handleMessage(update)
			}
		}
	}
}


func startMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Подсчет калорий", data: "calorie"},
		{name: "Тренировка", data: "traine"},
		{name: "Профиль", data: "profile"},
		{name: "Дневник", data: "dnevnik"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Меню «Тренировка»
func traineMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "🟢Тренировка: лёгкий уровень🟢", data: "Light"},
		{name: "🟡Тренировка: средний уровень🟡", data: "Midle"},
		{name: "🔴Тренировка: сложный уровень🔴", data: "Hard"},
		{name: "Назад", data: "back"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func enlightenment() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Спина (лёг.)", data: "backLight"},
		{name: "Руки (лёг.)", data: "handleLight"},
		{name: "Ноги (лёг.)", data: "kneesLight"},
		{name: "Трицепс (лёг.)", data: "tricepsLight"},
		{name: "Бицепс (лёг.)", data: "bicepsLight"},
		{name: "Назад", data: "lightBack"},
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func enlightenmentMidle() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Спина (ср.)", data: "backMidle"},
		{name: "Руки (ср.)", data: "handleMidle"},
		{name: "Ноги (ср.)", data: "kneesMidle"},
		{name: "Трицепс (ср.)", data: "tricepsMidle"},
		{name: "Бицепс (ср.)", data: "bicepsMidle"},
		{name: "Назад", data: "midleBack"},
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func enlightenmentHard() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Спина (сл.)", data: "backHard"},
		{name: "Руки (сл.)", data: "handleHard"},
		{name: "Ноги (сл.)", data: "kneesHard"},
		{name: "Трицепс (сл.)", data: "tricepsHard"},
		{name: "Бицепс (сл.)", data: "bicepsHard"},
		{name: "Назад", data: "hardBack"},
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func profileMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Заполнить профиль", data: "profile_anket"},
		{name: "Показать профиль", data: "Show_profile"},
		{name: "Назад", data: "back"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, st := range states {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(st.name, st.data),
		)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}


func handleCallback(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = bot.Send(del)

	switch data {
	case "calorie":
		prof, ok := userProfiles[chatID]
		if !ok || prof == nil || prof.Height == 0 || prof.Weight == 0 || prof.Age == 0 || prof.Gender == "" {
			sendText(chatID, "Ваш профиль ещё не заполнен. Сначала перейдите в «Профиль» и заполните данные.")
			return
		}

		genderForCalc := "женский"
		if strings.ToLower(prof.Gender) == "male" {
			genderForCalc = "мужской"
		}

		cals := calculateCalories(
			genderForCalc,
			float64(prof.Weight),
			float64(prof.Height),
			prof.Age,
		)
		msg := fmt.Sprintf("Ваш базовый обмен веществ: %.2f ккал/день\n(По данным профиля)", cals)
		sendText(chatID, msg)

	case "traine":
		msg := tgbotapi.NewMessage(chatID,
			"🏋️‍♂️ *Выберите уровень тренировки*\n"+
				"🟢 *Лёгкий уровень* — для начинающих и восстановления.\n"+
				"🟡 *Средний уровень* — для тех, кто готов к вызову.\n"+
				"🔴 *Сложный уровень* — для опытных и продвинутых.\n\n"+
				"🔙 *Назад* — вернуться в главное меню.")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "profile":
		msg := tgbotapi.NewMessage(chatID,
			"🧑‍💻 *Настройка профиля*\n"+
				"Выберите, что вы хотите сделать:\n\n"+
				"✏️ *Заполнить профиль* — внесите данные о своём росте, весе, возрасте и поле.\n"+
				"👀 *Показать профиль* — просмотрите ваши сохранённые данные.\n"+
				"🔙 *Назад* — вернуться в главное меню.")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = profileMenu()
		sendMessage(msg)

	case "Show_profile":
		prof, ok := userProfiles[chatID]
		if !ok || prof == nil {
			sendText(chatID, "Ваш профиль пока пуст. Нажмите «Заполнить профиль», чтобы внести данные.")
			return
		}
		msg := fmt.Sprintf("Ваш профиль:\nРост: %d см\nВес: %d кг\nВозраст: %d лет\nПол: %s",
			prof.Height, prof.Weight, prof.Age, prof.Gender)
		sendText(chatID, msg)

	case "back":
		mainMsg := tgbotapi.NewMessage(chatID,
			"👋 *Привет! Я ваш фитнес-помощник.*\n"+
				"Вот, что я умею:\n\n"+
				"🍎 *Подсчёт калорий* — помогу рассчитать дневную норму (по данным профиля).\n"+
				"🏋️‍♂️ *Тренировки* — подберу подходящие упражнения.\n"+
				"🧑‍💻 *Профиль* — сохраним ваши данные.\n"+
				"📔 *Дневник* — записывайте достижения и тренировки.\n\n"+
				"Выберите действие из меню или введите команду:\n"+
				"- /start — Главное меню\n"+
				"- /train — Перейти к тренировкам\n"+
				"- /profile — Настроить/посмотреть профиль\n"+
				"- /dnevnik — Записи и достижения.")
		mainMsg.ParseMode = "Markdown"
		mainMsg.ReplyMarkup = startMenu()
		sendMessage(mainMsg)

	case "profile_anket":
		startProfileWizard(chatID)

	case "dnevnik":
		msg := tgbotapi.NewMessage(chatID, "Вы находитесь в дневнике. Выберите действие:")
		msg.ReplyMarkup = diaryMenu()
		sendMessage(msg)

	case "diary_add":
		userStep[chatID] = diaryStepAdd
		sendText(chatID, "Введите, что вы сегодня сделали (например, какие упражнения выполнили).")

	case "diary_show":
		entries := userDiary[chatID]
		if len(entries) == 0 {
			sendText(chatID, "Пока нет записей в дневнике.")
			return
		}
		var sb strings.Builder
		sb.WriteString("Ваши записи в дневнике:\n\n")
		for i, entry := range entries {
			dateStr := entry.Date.Format("2006-01-02 15:04")
			sb.WriteString(fmt.Sprintf("%d) [%s] %s\n", i+1, dateStr, entry.Text))
		}
		sendText(chatID, sb.String())

		msg := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие в дневнике:")
		msg.ReplyMarkup = diaryMenu()
		sendMessage(msg)

	case "Light":
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали лёгкий уровень.")
		msg.ReplyMarkup = enlightenment()
		sendMessage(msg)

	case "backLight":
		msg := tgbotapi.NewMessage(chatID, "Пример лёгкой тренировки для спины:\n"+
			"1. «Кошка-корова» (Cat-Camel)\n2. Поза ребёнка (Child’s Pose)\n3. Superman (лёгкий вариант)\n"+
			"и т.д. ...")
		msg.ReplyMarkup = enlightenment()
		sendMessage(msg)

	case "handleLight":
		sendText(chatID, "Руки (лёг.) — аналогично, добавьте описание упражнений.")

	case "kneesLight":
		sendText(chatID, "Ноги (лёг.) — аналогично, добавьте описание упражнений.")

	case "tricepsLight":
		sendText(chatID, "Трицепс (лёг.) — аналогично.")

	case "bicepsLight":
		// Пример: можно читать текст из .env
		godotenv.Load()
		bigText := os.Getenv("BICEPSL")
		if bigText == "" {
			bigText = "Пустая переменная BICEPSL"
		}
		msg := tgbotapi.NewMessage(chatID, bigText)
		msg.ReplyMarkup = enlightenment()
		sendMessage(msg)

	case "lightBack":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "Midle":
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали средний уровень.")
		msg.ReplyMarkup = enlightenmentMidle()
		sendMessage(msg)

	case "backMidle":
		sendText(chatID, "Спина (ср.) — аналогично, описание упражнений тут.")

	case "handleMidle":
		sendText(chatID, "Руки (ср.) — аналогично.")

	case "kneesMidle":
		sendText(chatID, "Ноги (ср.) — аналогично.")

	case "tricepsMidle":
		sendText(chatID, "Трицепс (ср.) — аналогично.")

	case "bicepsMidle":
		sendText(chatID, "Бицепс (ср.) — аналогично.")

	case "midleBack":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "Hard":
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали сложный уровень.")
		msg.ReplyMarkup = enlightenmentHard()
		sendMessage(msg)

	case "backHard":
		sendText(chatID, "Спина (сл.) — аналогично.")

	case "handleHard":
		sendText(chatID, "Руки (сл.) — аналогично.")

	case "kneesHard":
		sendText(chatID, "Ноги (сл.) — аналогично.")

	case "tricepsHard":
		sendText(chatID, "Трицепс (сл.) — аналогично.")

	case "bicepsHard":
		sendText(chatID, "Бицепс (сл.) — аналогично.")

	case "hardBack":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)
	}
}

func startProfileWizard(chatID int64) {
	// Создаём (или очищаем) профиль для нового пользователя
	if userProfiles[chatID] == nil {
		userProfiles[chatID] = &user_profile{}
	}

	userStep[chatID] = profileStepHeight
	sendText(chatID, "Давайте заполним ваш профиль.\nВведите ваш рост (в см):")
}

func handleCommands(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID,
			"👋 *Привет! Я ваш фитнес-помощник.*\n"+
				"Вот, что я умею:\n"+
				"🍎 *Подсчёт калорий* — помогу рассчитать дневную норму (по данным профиля).\n"+
				"🏋️‍♂️ *Тренировки* — подберу подходящие упражнения.\n"+
				"🧑‍💻 *Профиль* — сохраняю ваши данные.\n"+
				"📔 *Дневник* — записывайте достижения.\n\n"+
				"Выберите действие в меню или введите:\n"+
				"- /train — Уровни тренировок\n"+
				"- /profile — Профиль\n"+
				"- /dnevnik — Дневник",
		)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = startMenu()
		sendMessage(msg)

	case "train":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "profile":
		msg := tgbotapi.NewMessage(chatID,
			"Настройка профиля:\nВыберите пункт ниже:")
		msg.ReplyMarkup = profileMenu()
		sendMessage(msg)

	default:
		sendText(chatID, "Неизвестная команда: "+update.Message.Command())
	}
}

func handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch userStep[chatID] {
	case profileStepHeight:
		height, err := strconv.Atoi(text)
		if err != nil || height <= 0 {
			sendText(chatID, "Пожалуйста, введите корректное число (рост в см).")
			return
		}
		userProfiles[chatID].Height = height
		time.Sleep(1 * time.Second)
		sendText(chatID, "📝Рост сохранён!\nТеперь введите ваш вес (в кг):")
		userStep[chatID] = profileStepWeight
		return

	case profileStepWeight:
		weight, err := strconv.Atoi(text)
		if err != nil || weight <= 0 {
			sendText(chatID, "Пожалуйста, введите корректный вес (число).")
			return
		}
		userProfiles[chatID].Weight = weight
		sendText(chatID, "📝Вес сохранён!\nТеперь введите ваш возраст (полных лет):")
		userStep[chatID] = profileStepAge
		return

	case profileStepAge:
		age, err := strconv.Atoi(text)
		if err != nil || age <= 0 {
			sendText(chatID, "Пожалуйста, введите корректный возраст (число).")
			return
		}
		userProfiles[chatID].Age = age
		sendText(chatID, "📝Возраст сохранён!\nТеперь введите ваш пол (male/female):")
		userStep[chatID] = profileStepGender
		return

	case profileStepGender:
		g := strings.ToLower(text)
		if g != "male" && g != "female" {
			sendText(chatID, "Пожалуйста, введите 'male' или 'female'.")
			return
		}
		userProfiles[chatID].Gender = g
		sendText(chatID, "Отлично, все данные заполнены!\nТеперь вы можете проверить /profile -> «Показать профиль» или рассчитать калории.")
		userStep[chatID] = profileStepDone
		return
	}

	if userStep[chatID] == diaryStepAdd {
		entry := DiaryEntry{
			Date: time.Now(),
			Text: text,
		}
		userDiary[chatID] = append(userDiary[chatID], entry)
		userStep[chatID] = 0

		sendText(chatID, "Запись добавлена в ваш дневник!")
		msg := tgbotapi.NewMessage(chatID, "Вы снова в дневнике. Выберите действие:")
		msg.ReplyMarkup = diaryMenu()
		sendMessage(msg)
		return
	}

	// Если не совпало ни с одним из «шагов» — отвечаем как на обычный текст
	sendText(chatID, "Я получил ваше сообщение: "+text)
}

func sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

func sendMessage(msg tgbotapi.Chattable) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
