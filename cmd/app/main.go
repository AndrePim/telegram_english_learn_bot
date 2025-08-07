package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	botHandlers "github.com/AndrePim/telegram_english_learn_bot"
	"github.com/AndrePim/telegram_english_learn_bot/internal/repository"
	"github.com/AndrePim/telegram_english_learn_bot/internal/service"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env файла (если он существует)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Получаем конфигурацию из переменных окружения
	config := getConfig()

	// Подключаемся к базе данных
	db, err := repository.NewDatabase(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	wordRepo := repository.NewWordRepository(db)

	// Инициализируем сервисы
	userService := service.NewUserService(userRepo)
	wordService := service.NewWordService(wordRepo)

	// Инициализируем обработчики бота
	handlers := botHandlers.NewBotHandlers(userService, wordService)

	// Создаем бота
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}

	b, err := bot.New(config.BotToken, opts...)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Регистрируем обработчики команд
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.HelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, handlers.AddHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/words", bot.MatchTypeExact, handlers.WordsHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/quiz", bot.MatchTypeExact, handlers.QuizHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/review", bot.MatchTypeExact, handlers.ReviewHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypePrefix, handlers.DeleteHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stats", bot.MatchTypeExact, handlers.StatsHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/image", bot.MatchTypePrefix, handlers.ImageHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, handlers.CallbackHandler)

	// Создаем контекст для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Bot started successfully!")

	// Запускаем бота
	b.Start(ctx)
}

type Config struct {
	BotToken   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func getConfig() *Config {
	return &Config{
		BotToken:   getEnv("BOT_TOKEN", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "english_bot_db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
