package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/olabanji12-ojo/church-backend/models"
	"github.com/sirupsen/logrus"
)

type EmbeddingService struct {
	apiKey string
	model  string
}

func NewEmbeddingService() *EmbeddingService {
	apiKey := os.Getenv("Hugging_Face_Key")
	if apiKey == "" {
		apiKey = os.Getenv("HF_API_KEY")
	}

	return &EmbeddingService{
		apiKey: apiKey,
		model:  "sentence-transformers/all-MiniLM-L6-v2",
	}
}

// GetEmbedding returns the 384-dimensional vector embedding for the given text.
func (es *EmbeddingService) GetEmbedding(text string) ([]float32, error) {
	if es.apiKey == "" {
		logrus.Warn("⚠️ Hugging_Face_Key is not set! Returning mock zero-vector for local development.")
		mockVector := make([]float32, 384)
		mockVector[0] = 0.5 // Set a dummy value so it's not a pure zero vector
		return mockVector, nil
	}

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", es.model)

	requestBody, err := json.Marshal(map[string]interface{}{
		"inputs": text,
		"options": map[string]interface{}{
			"wait_for_model": true,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+es.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Warnf("⚠️ HF API request failed or timed out fast: %v. Using fallback vector.", err)
		fallback := make([]float32, 384)
		for i := 0; i < len(text) && i < 384; i++ {
			fallback[i] = float32(text[i]) / 255.0
		}
		return fallback, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse interface{}
		json.NewDecoder(resp.Body).Decode(&errResponse)
		logrus.Errorf("❌ HF API returned error status %d: %v", resp.StatusCode, errResponse)
		return nil, fmt.Errorf("hugging face api error: status %d", resp.StatusCode)
	}

	var rawBytes bytes.Buffer
	_, err = rawBytes.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	var embedding []float32
	// Hugging Face can return a flat 1D array, a 2D array, or a 3D array depending on inputs.
	// We handle all shapes to ensure robust parsing.
	if err := json.Unmarshal(rawBytes.Bytes(), &embedding); err != nil {
		// Try 2D array: [[val1, val2, ...]]
		var nested2D [][]float32
		if err2 := json.Unmarshal(rawBytes.Bytes(), &nested2D); err2 == nil && len(nested2D) > 0 {
			embedding = nested2D[0]
		} else {
			// Try 3D array: [[[val1, val2, ...]]]
			var nested3D [][][]float32
			if err3 := json.Unmarshal(rawBytes.Bytes(), &nested3D); err3 == nil && len(nested3D) > 0 && len(nested3D[0]) > 0 {
				embedding = nested3D[0][0]
			} else {
				logrus.Errorf("❌ Failed to decode HF response: %s", rawBytes.String())
				return nil, fmt.Errorf("failed to decode embedding: %v", err)
			}
		}
	}

	if len(embedding) != 384 {
		logrus.Warnf("⚠️ Embedding dimension mismatch. Expected 384, got %d. Padding or trimming.", len(embedding))
		if len(embedding) < 384 {
			padded := make([]float32, 384)
			copy(padded, embedding)
			embedding = padded
		} else {
			embedding = embedding[:384]
		}
	}

	return embedding, nil
}

// GenerateUserText compiles the relevant user profile attributes into a single string for embedding.
func (es *EmbeddingService) GenerateUserText(user *models.User) string {
	bio := user.Bio
	if bio == "" {
		bio = "Seeking a faith-based partner."
	}

	return fmt.Sprintf(
		"First Name: %s. Gender: %s. Denomination: %s. Church frequency: %s. Prayer: %s. Bible: %s. Intention: %s. Bio: %s.",
		user.FirstName,
		user.Gender,
		user.Denomination,
		user.ChurchFreq,
		user.PrayerFreq,
		user.BibleFreq,
		user.Intention,
		bio,
	)
}

// GeneratePartnerPreferenceText compiles the preference details and written partner description into a single string.
func (es *EmbeddingService) GeneratePartnerPreferenceText(user *models.User) string {
	denomination := user.PreferredDenomination
	if denomination == "" || denomination == "any" {
		denomination = "Any"
	}
	churchFreq := user.PreferredChurchFreq
	if churchFreq == "" || churchFreq == "any" {
		churchFreq = "Any"
	}
	
	text := fmt.Sprintf(
		"Preferred Denomination: %s. Preferred Church Frequency: %s. Preferred Age Range: %d to %d.",
		denomination,
		churchFreq,
		user.MinAgePref,
		user.MaxAgePref,
	)

	if user.PartnerPrefText != "" {
		text += " Written Partner Description: " + user.PartnerPrefText
	} else {
		text += " Seeking a faith-oriented partner."
	}

	return text
}
