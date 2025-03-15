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

func (m *MockURLService) RedirectURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func setupMockService() (*MockURLService, *URLHandler) {
	mockService := new(MockURLService)
	handler := NewURLHandler(mockService)

	return mockService, handler
}

func TestShortenURL_Success(t *testing.T) {
	mockService, handler := setupMockService()

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
	mockService, handler := setupMockService()

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

func TestRedirect_Success(t *testing.T) {
	mockService, handler := setupMockService()

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	req.SetPathValue("shortcode", "abc123")
	rec := httptest.NewRecorder()

	mockService.On("RedirectURL", req.Context(), "abc123").Return("https://test.test", nil)

	handler.Redirect(rec, req)

	mockService.AssertExpectations(t)

	assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	assert.Equal(t, "https://test.test", rec.Header().Get("Location"))
}

func TestRedirect_NotFound(t *testing.T) {
	mockService, handler := setupMockService()

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	req.SetPathValue("shortcode", "abc123")
	rec := httptest.NewRecorder()

	mockService.On("RedirectURL", req.Context(), "abc123").Return("", errors.New("not found"))

	handler.Redirect(rec, req)

	mockService.AssertExpectations(t)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "not found", strings.TrimSpace(rec.Body.String()))
}
