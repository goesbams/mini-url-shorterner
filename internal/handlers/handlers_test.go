package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"mini-url-shortener/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Error(1)
}

func TestShortenURL_Success(t *testing.T) {
	mockService := new(MockURLService)
	handler := NewURLHandler(mockService)

	requestBody := `{"original_url":"http://test.test"}`

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	mockService.On("ShortenURL", req.Context(), "http://test.test").Return("abc123", nil)

	handler.ShortenURL(rec, req)
	mockService.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.URLResponse
	json.NewDecoder(rec.Body).Decode(&response)

	assert.NotEmpty(t, response.ShortCode)
	assert.Equal(t, "abc123", response.ShortCode)
}

func TestShortenURL_EmptyRequest(t *testing.T) {
	mockService := new(MockURLService)
	handler := NewURLHandler(mockService)

	requestBody := `{"original_url": ""}`

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	mockService.On("ShortenURL", req.Context(), "").Return("", errors.New("original url cannot be empty"))

	handler.ShortenURL(rec, req)
	mockService.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "original url cannot be empty")
}
