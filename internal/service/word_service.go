package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AndrePim/telegram_english_learn_bot/internal/repository"
)

type WordService struct {
	wordRepo *repository.WordRepository
}

func NewWordService(wordRepo *repository.WordRepository) *WordService {
	return &WordService{wordRepo: wordRepo}
}

// AddWord добавляет новое слово
func (s *WordService) AddWord(userID int64, word, translation, context string) error {
	// Проверяем, что слово и перевод не пустые
	if strings.TrimSpace(word) == "" || strings.TrimSpace(translation) == "" {
		return fmt.Errorf("word and translation cannot be empty")
	}

	newWord := &repository.Word{
		UserID:      userID,
		Word:        strings.TrimSpace(word),
		Translation: strings.TrimSpace(translation),
		Context:     strings.TrimSpace(context),
	}

	return s.wordRepo.SaveWord(newWord)
}

// GetUserWords получает все слова пользователя
func (s *WordService) GetUserWords(userID int64) ([]*repository.Word, error) {
	return s.wordRepo.GetUserWords(userID)
}

// GetWordsForReview получает слова для повторения
func (s *WordService) GetWordsForReview(userID int64) ([]*repository.Word, error) {
	return s.wordRepo.GetWordsForReview(userID)
}

// UpdateWordReview обновляет статус повторения слова
func (s *WordService) UpdateWordReview(wordID int, correct bool) error {
	return s.wordRepo.UpdateWordReview(wordID, correct)
}

// DeleteWord удаляет слово
func (s *WordService) DeleteWord(wordID int, userID int64) error {
	return s.wordRepo.DeleteWord(wordID, userID)
}

// QuizQuestion представляет вопрос для теста
type QuizQuestion struct {
	WordID     int
	Question   string
	Options    []string
	CorrectIdx int
}

// GenerateQuiz генерирует тест для пользователя
func (s *WordService) GenerateQuiz(userID int64) (*QuizQuestion, error) {
	log.Printf("Generating quiz for user %d", userID)
	words, err := s.wordRepo.GetUserWords(userID) // Используем GetUserWords
	if err != nil {
		log.Printf("Failed to get words for quiz: %v", err)
		return nil, fmt.Errorf("failed to get words for quiz: %w", err)
	}
	if len(words) < 4 {
		log.Printf("Not enough words for quiz: %d", len(words))
		return nil, fmt.Errorf("need at least 4 words to generate quiz")
	}

	// Выбираем случайное слово
	r := randv2.New(randv2.NewSource(time.Now().UnixNano()))
	targetIdx := r.Intn(len(words))
	targetWord := words[targetIdx]

	// Создаем варианты ответов
	options := make([]string, 4)
	correctIdx := r.Intn(4)
	options[correctIdx] = targetWord.Translation

	// Отслеживаем использованные индексы
	usedIndices := map[int]bool{targetIdx: true}
	optionIdx := 0

	// Заполняем остальные варианты
	for optionIdx < 4 {
		if optionIdx == correctIdx {
			optionIdx++
			continue
		}
		randIdx := r.Intn(len(words))
		if !usedIndices[randIdx] {
			options[optionIdx] = words[randIdx].Translation
			usedIndices[randIdx] = true
			optionIdx++
		}
	}

	// Логируем варианты для отладки
	log.Printf("Quiz options: %v, correctIdx: %d", options, correctIdx)

	return &QuizQuestion{
		WordID:     targetWord.ID,
		Question:   fmt.Sprintf("Как переводится слово: %s?", targetWord.Word),
		Options:    options,
		CorrectIdx: correctIdx,
	}, nil
}
