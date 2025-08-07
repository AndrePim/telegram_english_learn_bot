package service

import (
	"fmt"
	"math/rand"
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
	words, err := s.wordRepo.GetWordsForReview(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get words for quiz: %w", err)
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("no words available for quiz")
	}

	// Выбираем случайное слово
	rand.Seed(time.Now().UnixNano())
	targetWord := words[rand.Intn(len(words))]

	// Получаем все слова пользователя для создания вариантов ответов
	allWords, err := s.wordRepo.GetUserWords(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all words: %w", err)
	}

	if len(allWords) < 4 {
		return nil, fmt.Errorf("need at least 4 words to generate quiz")
	}

	// Создаем варианты ответов
	options := []string{targetWord.Translation}
	correctIdx := 0

	// Добавляем 3 неправильных варианта
	for len(options) < 4 {
		randomWord := allWords[rand.Intn(len(allWords))]
		if randomWord.ID != targetWord.ID {
			// Проверяем, что такого варианта еще нет
			exists := false
			for _, option := range options {
				if option == randomWord.Translation {
					exists = true
					break
				}
			}
			if !exists {
				options = append(options, randomWord.Translation)
			}
		}
	}

	// Перемешиваем варианты
	for i := len(options) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
		if i == correctIdx {
			correctIdx = j
		} else if j == correctIdx {
			correctIdx = i
		}
	}

	return &QuizQuestion{
		WordID:     targetWord.ID,
		Question:   fmt.Sprintf("Как переводится слово: *%s*?", targetWord.Word),
		Options:    options,
		CorrectIdx: correctIdx,
	}, nil
}
