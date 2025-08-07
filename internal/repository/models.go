package repository

import "time"

// User представляет пользователя бота
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

// Word представляет слово для изучения
type Word struct {
	ID          int       `json:"id"`
	UserID      int64     `json:"user_id"`
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
	Context     string    `json:"context"`
	CreatedAt   time.Time `json:"created_at"`
	LastReview  time.Time `json:"last_review"`
	NextReview  time.Time `json:"next_review"`
	Interval    int       `json:"interval"`   // Интервал в днях для повторения
	Difficulty  int       `json:"difficulty"` // Сложность слова (0-5)
}

// Quiz представляет тест
type Quiz struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	WordID    int       `json:"word_id"`
	Correct   bool      `json:"correct"`
	CreatedAt time.Time `json:"created_at"`
}
