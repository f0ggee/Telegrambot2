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

// ===================== СТРУКТУРЫ И ГЛОБАЛЬНЫЕ ПЕРЕМЕННЫЕ =====================

// Шаги пошагового опроса (профиля)
const (
	profileStepHeight = 1 // Вопрос о росте
	profileStepWeight = 2 // Вопрос о весе
	profileStepAge    = 3 // Вопрос о возрасте
	profileStepGender = 4 // Вопрос о поле
	profileStepDone   = 5 // Профиль заполнен
)

// Шаги пошагового опроса (калорий)
const (
	calorieStepWeight = 101
	calorieStepHeight = 102
	calorieStepAge    = 103
	calorieStepGender = 104
	calorieStepDone   = 105
)

// userStep хранит текущий шаг (либо для профиля, либо для калорий) для каждого пользователя
var userStep = make(map[int64]int)

// user_profile хранит данные о пользователе (из первого кода)
type user_profile struct {
	Height int
	Weight int
	Age    int
	Gender string
}

// userProfiles хранит профили по chatID (из первого кода)
var userProfiles = make(map[int64]*user_profile)

// Структура для описания «кнопок» (из первого кода)
type button struct {
	name string
	data string
}

var bot *tgbotapi.BotAPI

// ===================== ЛОГИКА РАСЧЁТА КАЛОРИЙ (из второго кода) =====================

func calculateCalories(gender string, weight float64, height float64, age int) float64 {
	if gender == "мужской" {
		return 88.36 + (13.4 * weight) + (4.8 * height) - (5.7 * float64(age))
	}
	return 447.6 + (9.2 * weight) + (3.1 * height) - (4.3 * float64(age))
}

// Для калорийного опроса будем хранить временно данные в отдельной карте (как во втором коде)
var calorieData = make(map[int64]map[string]string)

// ===================== MAIN =====================
func main() {
	// 1. Загружаем токен из .env (если нет, то можно прописать напрямую)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not loaded (если токен прописан вручную, то всё ок)")
	}

	botToken := os.Getenv("TG_BOT_API")
	// Если хотите — можно захардкодить токен вместо env
	// botToken := "7182429562:...ваш_токен..."

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

	// 3. Цикл обработки входящих сообщений/колбэков
	for update := range updates {
		if update.CallbackQuery != nil {
			// Обработаем нажатие на инлайн-кнопку
			handleCallback(update)
		} else if update.Message != nil {
			// Обработаем обычное сообщение
			if update.Message.IsCommand() {
				handleCommands(update)
			} else {
				handleMessage(update)
			}
		}
	}
}

// ===================== МЕНЮ (из первого кода) =====================

// Главное меню
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
		{name: "Тренировка: лёгкий уровень", data: "Light"},
		{name: "Тренировка: средний уровень", data: "Midle"},
		{name: "Тренировка: сложный уровень", data: "Hard"},
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

// Лёгкий уровень тренировок
func enlightenment() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Прокачка бицепса", data: "Bicepslight"},
		{name: "Прокачка передней части руки", data: "handle up"},
		{name: "Прокачка средней части руки", data: "handle middle"},
		{name: "Прокачка задней части руки", data: "handle behind"},
		{name: "Прокачка трицепса", data: "upgrade triceps"},
		{name: "Назад", data: "back2"},
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

// Средний уровень
func enlightenmentMidle() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Прокачка бицепса (ср.)", data: "BicepslightM"},
		{name: "Прокачка передней части руки (ср.)", data: "handle upM"},
		{name: "Прокачка средней части руки (ср.)", data: "handle middleM"},
		{name: "Прокачка задней части руки (ср.)", data: "handle behindM"},
		{name: "Прокачка трицепса (ср.)", data: "upgrade tricepsM"},
		{name: "Назад", data: "back3"},
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

// Меню «Профиль»
func profileMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{name: "Заполнить профиль", data: "profile_anket"}, // <-- запускаем пошаговый опрос
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

// ===================== ОБРАБОТКА INLINE-КНОПОК (из первого кода, с доработкой) =====================

func handleCallback(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	// Удаляем предыдущее сообщение (с кнопками)
	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = bot.Send(del)

	switch data {
	// Главное меню
	case "calorie":
		// Оставляем старый текст "Здесь будет подсчет калорий (пока не реализовано)."
		sendText(chatID, " ")
		// А сразу после — запускаем «второй» опрос (из второго кода):
		startCalorieWizard(chatID)

	case "traine":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "profile":
		msg := tgbotapi.NewMessage(chatID, "Настройка профиля:\nВыберите один из пунктов ниже:")
		msg.ReplyMarkup = profileMenu()
		sendMessage(msg)

	case "Show_profile":
		prof, ok := userProfiles[chatID]
		if !ok || prof == nil {
			sendText(chatID, "Ваш профиль пока пуст. Нажмите «Заполнить профиль», чтобы внести данные.")
			return
		}
		msg := "Ваш профиль:\n"
		msg += "Рост: " + strconv.Itoa(prof.Height) + "\n"
		msg += "Вес: " + strconv.Itoa(prof.Weight) + "\n"
		msg += "Возраст: " + strconv.Itoa(prof.Age) + "\n"
		msg += "Пол: " + prof.Gender + "\n"

		sendText(chatID, msg)

	case "back":
		// Главное меню
		mainMsg := tgbotapi.NewMessage(chatID, "Привет! Я ваш помощник-бот. Вот что я умею:\n- 📚 Подсчет калорий\n- 🏋️‍♂️ Тренировки\n- 🧑‍💻 Профиль\n\nВыберите действие в меню или введите /start, /train, /profile.")
		mainMsg.ReplyMarkup = startMenu()
		sendMessage(mainMsg)

	case "profile_anket":
		startProfileWizard(chatID)

	// --- Кнопки тренировки ---
	case "Light":
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали лёгкий уровень.")
		msg.ReplyMarkup = enlightenment()
		sendMessage(msg)

	case "Midle":
		msg := tgbotapi.NewMessage(chatID, "Вы выбрали средний уровень.")
		msg.ReplyMarkup = enlightenmentMidle()
		sendMessage(msg)

	case "Hard":
		sendText(chatID, "Вы выбрали сложный уровень.")

	case "back2":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "Bicepslight":
		sendText(chatID, "Тренировка бицепса (лёгкий уровень).")

	case "back3":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "handle upM":
		sendText(chatID, "Ср. уровень, прокачка передней части руки.")

	case "BicepslightM":
		sendText(chatID, "Ср. уровень, прокачка бицепса.")
	}
}

// ===================== ЗАПУСК ОПРОСА ДЛЯ ПРОФИЛЯ (из первого кода) =====================

func startProfileWizard(chatID int64) {
	// Создаём (или очищаем) профиль для нового пользователя
	if userProfiles[chatID] == nil {
		userProfiles[chatID] = &user_profile{}
	}

	// Ставим на первый шаг — ввод роста
	userStep[chatID] = profileStepHeight

	sendText(chatID, "Давайте заполним ваш профиль.\nВведите ваш рост (в см):")
}

// ===================== ОБРАБОТКА КОМАНД (/start, /profile, /train, ...) =====================

func handleCommands(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID, "Привет! Я ваш помощник-бот.\nВот что я умею:\n- 📚 Подсчет калорий\n- 🏋️‍♂️ Список тренировок\n- 🧑‍💻 Профиль\n\nВыберите меню или введите команды: /train, /profile.")
		msg.ReplyMarkup = startMenu()
		sendMessage(msg)

	case "train":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		sendMessage(msg)

	case "profile":
		msg := tgbotapi.NewMessage(chatID, "Настройка профиля:\nВыберите пункт ниже:")
		msg.ReplyMarkup = profileMenu()
		sendMessage(msg)

	default:
		sendText(chatID, "Неизвестная команда: "+update.Message.Command())
	}
}

// ===================== ОБРАБОТКА ОБЫЧНЫХ СООБЩЕНИЙ =====================

func handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// 1) Проверяем, не в режиме ли пошагового опроса (wizard) для ПРОФИЛЯ
	switch userStep[chatID] {
	case profileStepHeight:
		height, err := strconv.Atoi(text)
		if err != nil || height <= 0 {
			sendText(chatID, "Пожалуйста, введите корректное число (рост в см).")
			return
		}
		userProfiles[chatID].Height = height

		// Просто пример задержки и т.д.
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
		if text != "male" && text != "female" {
			sendText(chatID, "Пожалуйста, введите 'male' или 'female'.")
			return
		}
		userProfiles[chatID].Gender = text

		sendText(chatID, "Отлично, все данные заполнены!\nТеперь можете посмотреть профиль через /profile во вкладки  «Показать профиль».")
		userStep[chatID] = profileStepDone
		return
	}

	// 2) Если не в режиме опроса профиля, проверяем — не в опросе ли калорий
	if userStep[chatID] >= calorieStepWeight && userStep[chatID] <= calorieStepGender {
		handleCalorieWizard(chatID, text)
		return
	}

	// 3) Иначе — обычный текст
	sendText(chatID, "Я получил ваше сообщение: "+text)
}

// ===================== ЛОГИКА «ВТОРОГО» БОТА: Опрашиваем для расчёта калорий =====================

// startCalorieWizard — начинаем опрос по калориям
func startCalorieWizard(chatID int64) {
	// Создаём или обнуляем карту с ответами для данного пользователя
	calorieData[chatID] = map[string]string{}

	// Ставим шаг = calorieStepWeight
	userStep[chatID] = calorieStepWeight

	// Сообщение из второго кода (не изменяем текст!)
	sendText(chatID, "Привет! Я помогу рассчитать твоё дневное количество калорий. Введи свой вес в кг:")
}

// handleCalorieWizard — пошаговая логика (взята из второго кода)
func handleCalorieWizard(chatID int64, userMsg string) {
	data := calorieData[chatID]

	switch userStep[chatID] {
	case calorieStepWeight:
		weight, err := strconv.ParseFloat(userMsg, 64)
		if err != nil {
			sendText(chatID, "Пожалуйста, введи вес в числовом формате (например: 70.5):")
			return
		}
		data["weight"] = strconv.FormatFloat(weight, 'f', 1, 64)

		sendText(chatID, "Теперь введи свой рост в сантиметрах:")
		userStep[chatID] = calorieStepHeight
		return

	case calorieStepHeight:
		height, err := strconv.ParseFloat(userMsg, 64)
		if err != nil {
			sendText(chatID, "Пожалуйста, введи рост в числовом формате (например: 175):")
			return
		}
		data["height"] = strconv.FormatFloat(height, 'f', 1, 64)

		sendText(chatID, "Укажи свой возраст в годах:")
		userStep[chatID] = calorieStepAge
		return

	case calorieStepAge:
		age, err := strconv.Atoi(userMsg)
		if err != nil {
			sendText(chatID, "Пожалуйста, введи возраст в числовом формате (например: 25):")
			return
		}
		data["age"] = strconv.Itoa(age)

		sendText(chatID, "Теперь укажи свой пол (мужской или женский):")
		userStep[chatID] = calorieStepGender
		return

	case calorieStepGender:
		gender := strings.ToLower(strings.TrimSpace(userMsg))
		if gender != "мужской" && gender != "женский" {
			sendText(chatID, "Пожалуйста, укажи свой пол: мужской или женский.")
			return
		}
		data["gender"] = gender

		// Все данные собраны, рассчитываем калории
		weightVal, _ := strconv.ParseFloat(data["weight"], 64)
		heightVal, _ := strconv.ParseFloat(data["height"], 64)
		ageVal, _ := strconv.Atoi(data["age"])
		genderVal := data["gender"]

		calories := calculateCalories(genderVal, weightVal, heightVal, ageVal)
		result := fmt.Sprintf("Твой базовый обмен веществ (калории в день): %.2f ккал.", calories)

		sendText(chatID, result)

		// Сбрасываем данные
		delete(calorieData, chatID)
		userStep[chatID] = 0 // выходим из режима опроса
	}
}

// ===================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ =====================

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
