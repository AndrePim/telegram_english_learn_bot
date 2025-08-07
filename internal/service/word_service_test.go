package service

import (
	"testing"
)

func TestWordService_AddWord_EmptyWord(t *testing.T) {
	// Создаем мок репозитория (в реальном проекте лучше использовать интерфейсы)
	wordService := &WordService{}
	
	err := wordService.AddWord(123, "", "translation", "context")
	if err == nil {
		t.Error("Expected error for empty word, got nil")
	}
	
	if err.Error() != "word and translation cannot be empty" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestWordService_AddWord_EmptyTranslation(t *testing.T) {
	wordService := &WordService{}
	
	err := wordService.AddWord(123, "word", "", "context")
	if err == nil {
		t.Error("Expected error for empty translation, got nil")
	}
	
	if err.Error() != "word and translation cannot be empty" {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

