package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(database *Database) *WordRepository {
	return &WordRepository{db: database.db}
}

// SaveWord сохраняет новое слово
func (r *WordRepository) SaveWord(word *Word) error {
	query := `
		INSERT INTO words (user_id, word, translation, context, next_review)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(query, word.UserID, word.Word, word.Translation, word.Context, time.Now().AddDate(0, 0, 1)).
		Scan(&word.ID, &word.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to save word: %w", err)
	}
	
	return nil
}

// GetUserWords получает все слова пользователя
func (r *WordRepository) GetUserWords(userID int64) ([]*Word, error) {
	query := `
		SELECT id, user_id, word, translation, context, created_at, last_review, next_review, interval, difficulty
		FROM words WHERE user_id = $1 ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user words: %w", err)
	}
	defer rows.Close()
	
	var words []*Word
	for rows.Next() {
		word := &Word{}
		err := rows.Scan(
			&word.ID, &word.UserID, &word.Word, &word.Translation, &word.Context,
			&word.CreatedAt, &word.LastReview, &word.NextReview, &word.Interval, &word.Difficulty,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}
	
	return words, nil
}

// GetWordsForReview получает слова для повторения
func (r *WordRepository) GetWordsForReview(userID int64) ([]*Word, error) {
	query := `
		SELECT id, user_id, word, translation, context, created_at, last_review, next_review, interval, difficulty
		FROM words WHERE user_id = $1 AND next_review <= $2 ORDER BY next_review ASC LIMIT 10
	`
	
	rows, err := r.db.Query(query, userID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get words for review: %w", err)
	}
	defer rows.Close()
	
	var words []*Word
	for rows.Next() {
		word := &Word{}
		err := rows.Scan(
			&word.ID, &word.UserID, &word.Word, &word.Translation, &word.Context,
			&word.CreatedAt, &word.LastReview, &word.NextReview, &word.Interval, &word.Difficulty,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}
	
	return words, nil
}

// UpdateWordReview обновляет информацию о повторении слова
func (r *WordRepository) UpdateWordReview(wordID int, correct bool) error {
	// Получаем текущее слово
	var interval, difficulty int
	query := `SELECT interval, difficulty FROM words WHERE id = $1`
	err := r.db.QueryRow(query, wordID).Scan(&interval, &difficulty)
	if err != nil {
		return fmt.Errorf("failed to get word for update: %w", err)
	}
	
	// Алгоритм интервального повторения (упрощенный SM-2)
	if correct {
		if interval == 1 {
			interval = 6
		} else {
			interval = int(float64(interval) * 2.5)
		}
		if difficulty > 0 {
			difficulty--
		}
	} else {
		interval = 1
		if difficulty < 5 {
			difficulty++
		}
	}
	
	nextReview := time.Now().AddDate(0, 0, interval)
	
	updateQuery := `
		UPDATE words SET 
			last_review = $1, 
			next_review = $2, 
			interval = $3, 
			difficulty = $4 
		WHERE id = $5
	`
	
	_, err = r.db.Exec(updateQuery, time.Now(), nextReview, interval, difficulty, wordID)
	if err != nil {
		return fmt.Errorf("failed to update word review: %w", err)
	}
	
	return nil
}

// DeleteWord удаляет слово
func (r *WordRepository) DeleteWord(wordID int, userID int64) error {
	query := `DELETE FROM words WHERE id = $1 AND user_id = $2`
	
	result, err := r.db.Exec(query, wordID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete word: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("word not found or not owned by user")
	}
	
	return nil
}

