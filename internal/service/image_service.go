package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ImageService struct {
	apiKey string
	client *http.Client
}

func NewImageService() *ImageService {
	return &ImageService{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		client: &http.Client{},
	}
}

type ImageRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// GenerateImage генерирует изображение по описанию слова
func (s *ImageService) GenerateImage(word string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	prompt := fmt.Sprintf("A simple, clear illustration of: %s. Educational style, suitable for language learning.", word)

	reqBody := ImageRequest{
		Prompt: prompt,
		N:      1,
		Size:   "256x256",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var imageResp ImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&imageResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(imageResp.Data) == 0 {
		return "", fmt.Errorf("no images generated")
	}

	return imageResp.Data[0].URL, nil
}
