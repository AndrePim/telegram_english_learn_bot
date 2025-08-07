package service

import (
	"context"
	"log"
	"time"

	"github.com/AndrePim/telegram_english_learn_bot/internal/repository"
	"github.com/go-telegram/bot"
)

type SchedulerService struct {
	bot         *bot.Bot
	userRepo    *repository.UserRepository
	wordService *WordService
}

func NewSchedulerService(bot *bot.Bot, userRepo *repository.UserRepository, wordService *WordService) *SchedulerService {
	return &SchedulerService{
		bot:         bot,
		userRepo:    userRepo,
		wordService: wordService,
	}
}

// StartDailyReminders запускает ежедневные напоминания
func (s *SchedulerService) StartDailyReminders(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Запускаем сразу при старте (для тестирования можно изменить на меньший интервал)
	s.sendDailyReminders(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sendDailyReminders(ctx)
		}
	}
}

func (s *SchedulerService) sendDailyReminders(ctx context.Context) {
	log.Println("Sending daily reminders...")

	// Здесь должна быть логика получения всех пользователей
	// Для упрощения пропускаем эту часть, но в реальном проекте
	// нужно добавить метод GetAllUsers в UserRepository

	// Пример отправки напоминания (нужен userID)
	// s.sendReminderToUser(ctx, userID)
}

func (s *SchedulerService) sendReminderToUser(ctx context.Context, userID int64) {
	words, err := s.wordService.GetWordsForReview(userID)
	if err != nil {
		log.Printf("Failed to get words for user %d: %v", userID, err)
		return
	}

	if len(words) > 0 {
		message := "🔔 Напоминание! У вас есть слова для повторения. Используйте /review для просмотра."

		_, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   message,
		})

		if err != nil {
			log.Printf("Failed to send reminder to user %d: %v", userID, err)
		}
	}
}
