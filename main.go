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

// Структура для профиля
type user_profile struct {
	Height int
	Weight int
	Traine int
	Gender string
	Age    int
}

// Глобальные карты и переменные
var (
	bot          *tgbotapi.BotAPI
	userState    = make(map[int64]string)            // Отслеживает состояние (что бот ждёт от пользователя)
	userProfiles = make(map[int64]*user_profile)     // Профиль для каждого пользователя
	userDiary    = make(map[int64]map[string]string) // Дневник: для каждого chatID храним записи (дата -> запись)
	calorieData  = make(map[int64]map[string]string) // Временные данные для подсчёта калорий
)

// Функция расчёта калорий (пример)
func calculateCalories(gender string, weight float64, height float64, age int) float64 {
	if gender == "мужской" {
		return 88.36 + (13.4 * weight) + (4.8 * height) - (5.7 * float64(age))
	}
	// Если не «мужской», считаем как для «женский»
	return 447.6 + (9.2 * weight) + (3.1 * height) - (4.3 * float64(age))
}

// ===== Меню (inline-кнопки) =====

// Главное меню
func startMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("📋Подсчет калорий📋", "calorie"),
			tgbotapi.NewInlineKeyboardButtonData("💪Тренировка💪", "traine"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("📖Дневник📖", "diary"),
			tgbotapi.NewInlineKeyboardButtonData("👤Профиль👤", "profile"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Меню «Тренировка»
func traineMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Тренировка: лёгкий уровень", "Light"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Тренировка: средний уровень", "Midle"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Тренировка: сложный уровень", "Hard"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Под-меню лёгкого уровня
func enlightenment() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка бицепса", "Bicepslight"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка передней части руки", "handle up"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка средней части руки", "handle middle"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка задней части руки", "handle behind"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка трицепса", "upgrade triceps"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back2"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Под-меню среднего уровня
func enlightenmentMidle() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка бицепса (ср.)", "BicepslightM"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка передней части руки (ср.)", "handle upM"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка средней части руки (ср.)", "handle middleM"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка задней части руки (ср.)", "handle behindM"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка трицепса (ср.)", "upgrade tricepsM"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back3"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Под-меню сложного уровня
func enlightenmentHard() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка бицепса (сл.)", "BicepslightH"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка передней части руки (сл.)", "handle upH"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка средней части руки (сл.)", "handle middleH"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка задней части руки (сл.)", "handle behindH"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Прокачка трицепса (сл.)", "upgrade tricepsH"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back3"), // Возвращает к traineMenu()
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Меню «Профиль»
func profileMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Введите ваш рост", "Ask_height"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Ведите ваш возраст", "Ask_age"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Введите ваш пол", "Ask_gender"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Введите ваш вес", "Ask_weight"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Сколько занимаетесь (пример)", "Ask_traine"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Показать профиль", "Show_profile"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// Меню «Дневник»
func diaryMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Добавить запись", "add_entry"),
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть записи", "view_entries"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env not loaded (можно игнорировать, если токен прописан вручную)")
	}

	botToken := os.Getenv("TG_BOT_API")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации TG Bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Не удалось получить UpdatesChan: %v", err)
	}

	log.Println("Бот запущен...")

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

// --- Обработка команд (например, /start) ---
func handleCommands(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID,
			"Привет! Я ваш помощник-бот. Вот что я умею:\n"+
				"- Подсчет калорий\n"+
				"- Список тренировок\n"+
				"- Дневник\n"+
				"- Профиль.\n\n"+
				"Используйте меню ниже или введите команды:\n"+
				"/start, /train, /profile")
		msg.ReplyMarkup = startMenu()
		bot.Send(msg)

	case "train":
		trainMsg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		trainMsg.ReplyMarkup = traineMenu()
		bot.Send(trainMsg)

	case "profile":
		profMsg := tgbotapi.NewMessage(chatID, "Настройки профиля:")
		profMsg.ReplyMarkup = profileMenu()
		bot.Send(profMsg)

	default:
		sendText(chatID, "Неизвестная команда.")
	}
}

// --- Обработка нажатий на inline-кнопки (CallbackQuery) ---
func handleCallback(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	data := update.CallbackQuery.Data

	// Удалим старое сообщение (где кнопки)
	msgID := update.CallbackQuery.Message.MessageID
	delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
	_, _ = bot.Send(delMsg)

	switch data {
	// --- Дневник ---
	case "diary":
		msg := tgbotapi.NewMessage(chatID, "Вы в дневнике:")
		msg.ReplyMarkup = diaryMenu()
		bot.Send(msg)

	case "add_entry":
		userState[chatID] = "adding_entry"
		sendText(chatID, "Напишите, что хотите записать (например: тренировка, съедено ...).")

	case "view_entries":
		entries, exists := userDiary[chatID]
		if !exists || len(entries) == 0 {
			sendText(chatID, "Ваш дневник пока пуст.")
			return
		}
		var response string
		response = "Ваши записи:\n"
		for date, entry := range entries {
			response += fmt.Sprintf("%s: %s\n", date, entry)
		}
		sendText(chatID, response)

	// --- Подсчёт калорий ---
	case "calorie":
		calorieData[chatID] = make(map[string]string)
		userState[chatID] = "calorie_weight"
		sendText(chatID, "Введите свой вес (кг), например 70.5:")

	// --- ТРЕНИРОВКИ (выбор уровня) ---
	case "traine":
		msg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		msg.ReplyMarkup = traineMenu()
		bot.Send(msg)

	case "Light":
		lightMsg := tgbotapi.NewMessage(chatID, "Вы выбрали лёгкий уровень.")
		lightMsg.ReplyMarkup = enlightenment()
		bot.Send(lightMsg)

	case "Midle":
		midMsg := tgbotapi.NewMessage(chatID, "Вы выбрали средний уровень.")
		midMsg.ReplyMarkup = enlightenmentMidle()
		bot.Send(midMsg)

	case "Hard":
		// Показываем меню сложного уровня
		hardMsg := tgbotapi.NewMessage(chatID, "Вы выбрали сложный уровень.")
		hardMsg.ReplyMarkup = enlightenmentHard()
		bot.Send(hardMsg)

	// --- Назад из под-меню лёгкого уровня ---
	case "back2":
		backMsg := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		backMsg.ReplyMarkup = traineMenu()
		bot.Send(backMsg)

	// --- Назад из под-меню среднего/сложного уровня ---
	case "back3":
		backMsg2 := tgbotapi.NewMessage(chatID, "Выберите уровень тренировки:")
		backMsg2.ReplyMarkup = traineMenu()
		bot.Send(backMsg2)

	// --- Лёгкий уровень, 5 мышечных групп ---
	case "Bicepslight":
		textBicepsLight := `
Лёгкий уровень (Новичок) для бицепса:
1. Сгибания рук с гантелями стоя 
   - 2–3 подхода по 12–15 повторений
2. Концентрированные сгибания (поочерёдно)
   - 2 подхода по 10–12 повторений
3. Подтягивания обратным хватом (если можете)
   - 2 подхода по макс. (5–8 для начала)

Следите за техникой и не берите большой вес!
`
		sendText(chatID, textBicepsLight)

	case "handle up":
		textFrontArmLight := `
Лёгкий уровень: передняя часть руки (передние дельты):
1. Жим гантелей сидя (или стоя)
   - 2–3 подхода по 12–15 повторений
2. Передние подъёмы гантелей
   - 2 подхода по 10–12 повторений
3. Отжимания от пола с узкой постановкой
   - 2 подхода по 10–12

Не спешите, главное — техника!
`
		sendText(chatID, textFrontArmLight)

	case "handle middle":
		textMiddleDeltsLight := `
Лёгкий уровень: средняя часть руки (средние дельты):
1. Разведения рук в стороны (гантели) стоя
   - 2–3 подхода по 12–15 повторений
2. Упрощённые "Арчер отжимания"
   - 2 подхода по 5–8 на каждую сторону
3. Подъём гантелей через стороны в наклоне
   - 2 подхода по 10–12

Разогревайте плечи перед нагрузкой.
`
		sendText(chatID, textMiddleDeltsLight)

	case "handle behind":
		textBehindDeltsLight := `
Лёгкий уровень: задняя часть руки (задняя дельта):
1. Разведения в наклоне с гантелями
   - 2–3 подхода по 12–15
2. Тяга гантели к поясу узким хватом
   - 2 подхода по 10–12
3. Обратные отжимания от скамьи (упрощённые)
   - 2 подхода по 10

Избегайте рывков, держите спину ровно.
`
		sendText(chatID, textBehindDeltsLight)

	case "upgrade triceps":
		textTricepsLight := `
Лёгкий уровень для трицепса:
1. Отжимания узким хватом 
   - 2–3 подхода по 10–12
2. Французский жим с гантелью (одной рукой)
   - 2 подхода по 10–12
3. Обратные отжимания от скамьи (ноги на полу)
   - 2 подхода по 10

Держите локти ближе к туловищу.
`
		sendText(chatID, textTricepsLight)

	// --- Средний уровень, 5 мышечных групп ---
	case "BicepslightM":
		textBicepsMiddle := `
Средний уровень для бицепса:
1. Сгибания рук со штангой 
   - 3 подхода по 8–10
2. Сгибания "Молоток" (гантели)
   - 3 подхода по 10–12
3. Подтягивания обратным хватом
   - 3 подхода по 8–10

Следите за разминкой.
`
		sendText(chatID, textBicepsMiddle)

	case "handle upM":
		textFrontArmMiddle := `
Средний уровень: передняя часть руки (передние дельты):
1. Армейский жим штанги стоя
   - 3 подхода по 8–10
2. Передние подъёмы гантелей (наклонная)
   - 3 подхода по 10–12
3. Отжимания в стойке у стены (упрощённо)
   - 3 подхода по 8–10

Увеличивайте вес постепенно.
`
		sendText(chatID, textFrontArmMiddle)

	case "handle middleM":
		textMiddleDeltsMiddle := `
Средний уровень: средняя часть руки (средние дельты):
1. Разведения рук в стороны (средний вес)
   - 3 подхода по 10–12
2. Разведения в тренажёре "бабочка" (дельты)
   - 3 подхода по 10–12
3. Жим Арнольда (гантели)
   - 3 подхода по 8–10
`
		sendText(chatID, textMiddleDeltsMiddle)

	case "handle behindM":
		textBehindDeltsMiddle := `
Средний уровень: задняя дельта:
1. Разведения в наклоне с гантелями (средний вес)
   - 3 подхода по 10–12
2. Тяга штанги в наклоне узким хватом
   - 3 подхода по 8–10
3. Горизонтальные подтягивания обратным хватом
   - 3 подхода по 8–12
`
		sendText(chatID, textBehindDeltsMiddle)

	case "upgrade tricepsM":
		textTricepsMiddle := `
Средний уровень для трицепса:
1. Жим штанги узким хватом
   - 3 подхода по 8–10
2. Французский жим EZ-штангой (лёжа)
   - 3 подхода по 8–10
3. Обратные отжимания от скамьи (ноги повыше)
   - 3 подхода по 12
`
		sendText(chatID, textTricepsMiddle)

	// --- Сложный уровень, 5 мышечных групп ---
	case "BicepslightH":
		textBicepsHard := `
Сложный уровень для бицепса:
1. Сгибания рук со штангой на наклонной скамье
   - 4 подхода по 6–8
2. Подтягивания обратным хватом с отягощением
   - 4 подхода по 8–10
3. "21" на бицепс (7 нижних, 7 верхних, 7 полных)
   - 2 подхода по 21
`
		sendText(chatID, textBicepsHard)

	case "handle upH":
		textFrontArmHard := `
Сложный уровень: передняя часть руки (передние дельты):
1. Жим штанги стоя (heavy)
   - 4 подхода по 6–8
2. Передние махи со штангой
   - 4 подхода по 8–10
3. Отжимания в стойке на руках (Handstand push-ups)
   - 3 подхода по макс.
`
		sendText(chatID, textFrontArmHard)

	case "handle middleH":
		textMiddleHard := `
Сложный уровень: средние дельты:
1. Разведения рук с тяжёлыми гантелями
   - 4 подхода по 6–8
2. Подъём гантелей через стороны в блоке
   - 3 подхода по 10–12
3. Статическое удержание (в стороны) 20–30 сек
   - 2–3 подхода
`
		sendText(chatID, textMiddleHard)

	case "handle behindH":
		textBehindHard := `
Сложный уровень: задняя дельта:
1. Разведения в наклоне с тяжёлыми гантелями
   - 4 подхода по 6–8
2. Обратные "бабочки" в тренажёре
   - 4 подхода по 8–10
3. Подтягивания широким хватом за голову
   - 3 подхода по 8–10 (осторожно с техникой)
`
		sendText(chatID, textBehindHard)

	case "upgrade tricepsH":
		textTricepsHard := `
Сложный уровень для трицепса:
1. Жим штанги узким хватом (heavy)
   - 4 подхода по 6–8
2. Французский жим стоя (тяжёлая гантель)
   - 4 подхода по 8
3. Отжимания на брусьях с отягощением
   - 4 подхода по 8–10
`
		sendText(chatID, textTricepsHard)

	// --- Профиль ---
	case "profile":
		msg := tgbotapi.NewMessage(chatID, "Работа с профилем:")
		msg.ReplyMarkup = profileMenu()
		bot.Send(msg)

	case "Ask_height":
		userState[chatID] = "asking_height"
		sendText(chatID, "Введите ваш рост (например: 170).")

	case "Ask_age":
		userState[chatID] = "asking_age"
		sendText(chatID, "Введите ваш возраст (например: 25).")

	case "Ask_gender":
		userState[chatID] = "asking_gender"
		sendText(chatID, "Введите ваш пол (мужской/женский).")

	case "Ask_weight":
		userState[chatID] = "asking_weight"
		sendText(chatID, "Введите ваш вес (например: 70).")

	case "Ask_traine":
		userState[chatID] = "asking_traine"
		sendText(chatID, "Сколько лет вы занимаетесь в зале? (пример: 3 года).")

	case "Show_profile":
		prof, ok := userProfiles[chatID]
		if !ok {
			sendText(chatID, "Ваш профиль пока пуст.")
			return
		}
		message := "Ваш профиль:\n"
		message += fmt.Sprintf("Рост: %d \n", prof.Height)
		message += fmt.Sprintf("Вес: %d \n", prof.Weight)
		message += fmt.Sprintf("Возраст: %d \n", prof.Age)
		message += fmt.Sprintf("Пол: %s \n", prof.Gender)
		message += fmt.Sprintf("Тренировки: %d\n", prof.Traine)
		sendText(chatID, message)

	// --- Назад в главное меню ---
	case "back":
		mainMsg := tgbotapi.NewMessage(chatID,
			"Привет! Я ваш помощник-бот. Вот что я умею:\n- 📚 Подсчет калорий.\n- 🏋️‍♂️ Список тренировок.\n- 🧑‍💻 Профиль.\n- 📙 Дневник.\n\nИспользуйте команду /start, чтобы начать заново.")
		mainMsg.ReplyMarkup = startMenu()
		bot.Send(mainMsg)
	}
}

// --- Обработка обычных сообщений (не команд, не колбэков) ---
func handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch userState[chatID] {
	// --- Дневник ---
	case "adding_entry":
		if userDiary[chatID] == nil {
			userDiary[chatID] = make(map[string]string)
		}
		date := time.Now().Format("02-01-2006") // "дд-мм-гггг"
		userDiary[chatID][date] = text

		sendText(chatID, "Запись добавлена в дневник!")
		userState[chatID] = ""

	// --- Подсчёт калорий (шаги) ---
	case "calorie_weight":
		weight, err := strconv.ParseFloat(text, 64)
		if err != nil {
			sendText(chatID, "Ошибка: введите число (например: 70.5).")
			return
		}
		calorieData[chatID]["weight"] = fmt.Sprintf("%.1f", weight)

		userState[chatID] = "calorie_height"
		sendText(chatID, "Введите свой рост в см (например: 175):")
		return

	case "calorie_height":
		height, err := strconv.ParseFloat(text, 64)
		if err != nil {
			sendText(chatID, "Ошибка: введите число (например: 175).")
			return
		}
		calorieData[chatID]["height"] = fmt.Sprintf("%.1f", height)

		userState[chatID] = "calorie_age"
		sendText(chatID, "Введите свой возраст (например: 25):")
		return

	case "calorie_age":
		age, err := strconv.Atoi(text)
		if err != nil {
			sendText(chatID, "Ошибка: введите целое число (например: 25).")
			return
		}
		calorieData[chatID]["age"] = strconv.Itoa(age)

		userState[chatID] = "calorie_gender"
		sendText(chatID, "Укажите свой пол (мужской/женский):")
		return

	case "calorie_gender":
		gender := strings.ToLower(strings.TrimSpace(text))
		if gender != "мужской" && gender != "женский" {
			sendText(chatID, "Пожалуйста, укажите 'мужской' или 'женский'.")
			return
		}
		calorieData[chatID]["gender"] = gender

		// Все данные собраны
		w, _ := strconv.ParseFloat(calorieData[chatID]["weight"], 64)
		h, _ := strconv.ParseFloat(calorieData[chatID]["height"], 64)
		a, _ := strconv.Atoi(calorieData[chatID]["age"])

		res := calculateCalories(gender, w, h, a)
		sendText(chatID, fmt.Sprintf("Ваш базовый обмен веществ: %.2f ккал в день.", res))

		// Очищаем данные и сбрасываем состояние
		delete(calorieData, chatID)
		userState[chatID] = ""

	// --- Профиль (ввод роста, веса, возраста, пола, стажа) ---
	case "asking_height":
		h, err := strconv.Atoi(text)
		if err != nil {
			sendText(chatID, "Введите целое число (например: 170).")
			return
		}
		if userProfiles[chatID] == nil {
			userProfiles[chatID] = &user_profile{}
		}
		userProfiles[chatID].Height = h

		sendText(chatID, "Ваш рост сохранён!")
		userState[chatID] = ""

	case "asking_age":
		a, err := strconv.Atoi(text)
		if err != nil {
			sendText(chatID, "Введите целое число (например: 25).")
			return
		}
		if userProfiles[chatID] == nil {
			userProfiles[chatID] = &user_profile{}
		}
		userProfiles[chatID].Age = a

		sendText(chatID, "Ваш возраст сохранён!")
		userState[chatID] = ""

	case "asking_gender":
		g := strings.ToLower(strings.TrimSpace(text))
		if g != "мужской" && g != "женский" {
			sendText(chatID, "Укажите 'мужской' или 'женский'.")
			return
		}
		if userProfiles[chatID] == nil {
			userProfiles[chatID] = &user_profile{}
		}
		userProfiles[chatID].Gender = g

		sendText(chatID, "Ваш пол сохранён!")
		userState[chatID] = ""

	case "asking_weight":
		w, err := strconv.Atoi(text)
		if err != nil {
			sendText(chatID, "Введите целое число (например: 70).")
			return
		}
		if userProfiles[chatID] == nil {
			userProfiles[chatID] = &user_profile{}
		}
		userProfiles[chatID].Weight = w

		sendText(chatID, "Ваш вес сохранён!")
		userState[chatID] = ""

	case "asking_traine":
		tr, err := strconv.Atoi(text)
		if err != nil {
			sendText(chatID, "Введите целое число (например: 3).")
			return
		}
		if tr <= 3 {
			sendText(chatID, "Вы большой молодец, что начали заниматься! Рекомендуем лёгкий уровень.")
		} else if tr <= 5 {
			sendText(chatID, "Отлично, вы продолжаете заниматься! Средний уровень — для вас.")
		} else {
			sendText(chatID, "Вы уже хорошо подготовлены, можно пробовать сложный уровень!")
		}

		if userProfiles[chatID] == nil {
			userProfiles[chatID] = &user_profile{}
		}
		userProfiles[chatID].Traine = tr
		sendText(chatID, "Данные о тренировках сохранены!")
		userState[chatID] = ""

	default:
		// Если бот не в режиме ввода, просто отвечаем
		sendText(chatID, "Я получил ваше сообщение: "+text)
	}
}

// Вспомогательная функция для отправки текста
func sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
