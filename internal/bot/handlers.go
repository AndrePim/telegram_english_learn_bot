package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/AndrePim/telegram_english_learn_bot/internal/repository"
	"github.com/AndrePim/telegram_english_learn_bot/internal/service"
)

type BotHandlers struct {
	userService *service.UserService
	wordService *service.WordService
}

func NewBotHandlers(userService *service.UserService, wordService *service.WordService) *BotHandlers {
	return &BotHandlers{
		userService: userService,
		wordService: wordService,
	}
}

// DefaultHandler обрабатывает неизвестные команды
func (h *BotHandlers) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Извините, я не понимаю эту команду. Используйте /help для получения справки.",
	})
}

// StartHandler обрабатывает команду /start
func (h *BotHandlers) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From

	// Регистрируем пользователя
	err := h.userService.RegisterUser(user.ID, user.Username, user.FirstName, user.LastName)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Произошла ошибка при регистрации. Попробуйте позже.",
		})
		return
	}

	welcomeText := fmt.Sprintf(`Привет, %s! 👋

Я бот для изучения английского языка. Вот что я умею:

📝 /add - Добавить новое слово
📚 /words - Посмотреть все ваши слова
🧠 /quiz - Пройти тест
🔄 /review - Повторить слова
❓ /help - Показать справку

Начните с добавления слов командой /add!`, user.FirstName)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
	})
}

// HelpHandler обрабатывает команду /help
func (h *BotHandlers) HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	helpText := `🤖 Справка по боту для изучения английского

📝 /add - Добавить новое слово
   Формат: /add слово - перевод
   Пример: /add apple - яблоко

📚 /words - Показать все ваши слова

🧠 /quiz - Пройти тест на знание слов

🔄 /review - Повторить слова, которые пора повторить

🗑️ /delete [номер] - Удалить слово по номеру из списка

📊 /stats - Показать статистику изучения

🎨 /image [слово] - Сгенерировать изображение для слова

❓ /help - Показать эту справку

💡 Совет: Добавляйте контекст к словам для лучшего запоминания!
Пример: /add beautiful - красивый (She is beautiful)`

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
}

// AddHandler обрабатывает команду /add
func (h *BotHandlers) AddHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/add"))

	if text == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Используйте формат: /add слово - перевод\nПример: /add apple - яблоко",
		})
		return
	}

	// Парсим слово и перевод
	parts := strings.Split(text, " - ")
	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Используйте формат: /add слово - перевод\nПример: /add apple - яблоко",
		})
		return
	}

	word := strings.TrimSpace(parts[0])
	translation := strings.TrimSpace(parts[1])
	context := ""

	// Если есть контекст в скобках
	if len(parts) > 2 {
		context = strings.TrimSpace(strings.Join(parts[2:], " - "))
	}

	err := h.wordService.AddWord(userID, word, translation, context)
	if err != nil {
		log.Printf("Failed to add word: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при добавлении слова. Попробуйте еще раз.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ Слово '%s' добавлено!", word),
	})
}

// WordsHandler обрабатывает команду /words
func (h *BotHandlers) WordsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	words, err := h.wordService.GetUserWords(userID)
	if err != nil {
		log.Printf("Failed to get user words: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при получении слов.",
		})
		return
	}

	if len(words) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "У вас пока нет сохраненных слов. Добавьте их командой /add!",
		})
		return
	}

	var response strings.Builder
	response.WriteString("📚 Ваши слова:\n\n")

	for i, word := range words {
		response.WriteString(fmt.Sprintf("%d. *%s* - %s", i+1, word.Word, word.Translation))
		if word.Context != "" {
			response.WriteString(fmt.Sprintf(" (%s)", word.Context))
		}
		response.WriteString("\n")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response.String(),
		ParseMode: models.ParseModeMarkdown,
	})
}

// QuizHandler обрабатывает команду /quiz
func (h *BotHandlers) QuizHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	quiz, err := h.wordService.GenerateQuiz(userID)
	if err != nil {
		log.Printf("Failed to generate quiz: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось создать тест. Убедитесь, что у вас есть минимум 4 слова.",
		})
		return
	}

	// Создаем inline клавиатуру с вариантами ответов
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, len(quiz.Options)),
	}

	for i, option := range quiz.Options {
		keyboard.InlineKeyboard[i] = []models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d. %s", i+1, option),
				CallbackData: fmt.Sprintf("quiz_%d_%d_%d", quiz.WordID, i, quiz.CorrectIdx),
			},
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        quiz.Question,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: keyboard,
	})
}

// ReviewHandler обрабатывает команду /review
func (h *BotHandlers) ReviewHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	words, err := h.wordService.GetWordsForReview(userID)
	if err != nil {
		log.Printf("Failed to get words for review: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при получении слов для повторения.",
		})
		return
	}

	if len(words) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🎉 Отлично! Сейчас нет слов для повторения. Проверьте позже или добавьте новые слова!",
		})
		return
	}

	var response strings.Builder
	response.WriteString("🔄 Слова для повторения:\n\n")

	for i, word := range words {
		response.WriteString(fmt.Sprintf("%d. *%s* - %s", i+1, word.Word, word.Translation))
		if word.Context != "" {
			response.WriteString(fmt.Sprintf(" (%s)", word.Context))
		}
		response.WriteString("\n")
	}

	response.WriteString("\n💡 Пройдите тест командой /quiz для закрепления!")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      response.String(),
		ParseMode: models.ParseModeMarkdown,
	})
}

// CallbackHandler обрабатывает callback запросы (ответы на тесты)
func (h *BotHandlers) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	data := callback.Data

	// Парсим данные callback'а
	parts := strings.Split(data, "_")
	if len(parts) != 4 || parts[0] != "quiz" {
		return
	}

	wordID, _ := strconv.Atoi(parts[1])
	selectedIdx, _ := strconv.Atoi(parts[2])
	correctIdx, _ := strconv.Atoi(parts[3])

	correct := selectedIdx == correctIdx

	// Обновляем статистику слова
	err := h.wordService.UpdateWordReview(wordID, correct)
	if err != nil {
		log.Printf("Failed to update word review: %v", err)
	}

	var responseText string
	if correct {
		responseText = "✅ Правильно! Отличная работа!"
	} else {
		responseText = "❌ Неправильно. Не расстраивайтесь, продолжайте изучать!"
	}

	// Отвечаем на callback
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
		Text:            responseText,
		ShowAlert:       true,
	})

	// Обновляем сообщение
	if msg := callback.Message.Message; msg != nil {
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
			Text:      fmt.Sprintf("%s\n\n%s", msg.Text, responseText),
			ParseMode: models.ParseModeMarkdown,
		})
	}
}

// DeleteHandler обрабатывает команду /delete
func (h *BotHandlers) DeleteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/delete"))

	if text == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Используйте формат: /delete [номер]\nПример: /delete 1\n\nДля просмотра номеров слов используйте /words",
		})
		return
	}

	wordNum, err := strconv.Atoi(text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Неверный номер слова. Используйте /words для просмотра списка.",
		})
		return
	}

	// Получаем слова пользователя
	words, err := h.wordService.GetUserWords(userID)
	if err != nil {
		log.Printf("Failed to get user words: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при получении слов.",
		})
		return
	}

	if wordNum < 1 || wordNum > len(words) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Неверный номер. У вас %d слов. Используйте /words для просмотра.", len(words)),
		})
		return
	}

	wordToDelete := words[wordNum-1]
	err = h.wordService.DeleteWord(wordToDelete.ID, userID)
	if err != nil {
		log.Printf("Failed to delete word: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при удалении слова.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ Слово '%s' удалено!", wordToDelete.Word),
	})
}

// StatsHandler обрабатывает команду /stats
func (h *BotHandlers) StatsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID

	words, err := h.wordService.GetUserWords(userID)
	if err != nil {
		log.Printf("Failed to get user words: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при получении статистики.",
		})
		return
	}

	reviewWords, err := h.wordService.GetWordsForReview(userID)
	if err != nil {
		log.Printf("Failed to get words for review: %v", err)
		reviewWords = []*repository.Word{} // Продолжаем с пустым списком
	}

	totalWords := len(words)
	wordsForReview := len(reviewWords)

	var response strings.Builder
	response.WriteString("📊 Ваша статистика:\n\n")
	response.WriteString(fmt.Sprintf("📚 Всего слов: %d\n", totalWords))
	response.WriteString(fmt.Sprintf("🔄 Слов для повторения: %d\n", wordsForReview))
	response.WriteString(fmt.Sprintf("✅ Изученных слов: %d\n", totalWords-wordsForReview))

	if totalWords > 0 {
		progress := float64(totalWords-wordsForReview) / float64(totalWords) * 100
		response.WriteString(fmt.Sprintf("📈 Прогресс: %.1f%%\n", progress))
	}

	response.WriteString("\n💡 Продолжайте изучать новые слова!")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   response.String(),
	})
}

// ImageHandler обрабатывает команду /image
func (h *BotHandlers) ImageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/image"))

	if text == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Используйте формат: /image [слово]\nПример: /image apple",
		})
		return
	}

	// Отправляем сообщение о том, что генерируем изображение
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "🎨 Генерирую изображение... Это может занять несколько секунд.",
	})

	// Здесь должна быть интеграция с ImageService
	// Для демонстрации отправляем заглушку
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("🖼️ Изображение для слова '%s' будет здесь!\n\n💡 Для активации этой функции настройте OPENAI_API_KEY в переменных окружения.", text),
	})
}
