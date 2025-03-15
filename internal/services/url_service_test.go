package services

import (
	"context"
	"errors"
	"mini-url-shortener/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLRepository struct {
	mock.Mock
}

func (r *MockURLRepository) Begin(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}

func (r *MockURLRepository) Commit() error {
	args := r.Called()
	return args.Error(0)
}

func (r *MockURLRepository) Rollback() error {
	args := r.Called()
	return args.Error(0)
}

func (m *MockURLRepository) CreateShortCode(ctx context.Context, url *models.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) FindByShortCode(ctx context.Context, shortCode string) (url *models.URL, err error) {
	args := m.Called(ctx, shortCode)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *MockURLRepository) UpdateClickByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func setupMockRepository() (*MockURLRepository, URLService) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

	return mockRepo, service
}

func TestShortenURL_Success(t *testing.T) {
	mockRepo, service := setupMockRepository()

	url := &models.URL{
		OriginalURL: "https://test.test",
		ShortCode:   "a3f6fe",
	}

	mockRepo.On("CreateShortCode", context.Background(), url).Return(nil)

	shortCode, err := service.ShortenURL(context.Background(), url.OriginalURL)
	assert.NoError(t, err)

	assert.Equal(t, url.ShortCode, shortCode)
	assert.Equal(t, "a3f6fe", shortCode)
	mockRepo.AssertExpectations(t)
}

func TestShortenURL_EmptyURL(t *testing.T) {
	mockRepo, service := setupMockRepository()

	_, err := service.ShortenURL(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, "original url cannot be empty", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestShortenURL_RepositoryError(t *testing.T) {
	mockRepo, service := setupMockRepository()

	url := &models.URL{
		OriginalURL: "https://test.test",
		ShortCode:   "a3f6fe",
	}

	mockRepo.On("CreateShortCode", context.Background(), url).Return(errors.New("database error"))

	shortCode, err := service.ShortenURL(context.Background(), url.OriginalURL)
	assert.Error(t, err)

	assert.Empty(t, shortCode)
	assert.Equal(t, "error creating shorten url: database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRedirectURL_Success(t *testing.T) {
	mockRepo, service := setupMockRepository()
	ctx := context.Background()

	url := &models.URL{
		ID:          1,
		OriginalURL: "https://test.test",
		ShortCode:   "abc123",
	}

	mockRepo.On("Begin", ctx).Return(nil)
	mockRepo.On("FindByShortCode", ctx, url.ShortCode).Return(url, nil)
	mockRepo.On("UpdateClickByID", ctx, url.ID).Return(nil)
	mockRepo.On("Commit").Return(nil)

	originalURL, err := service.RedirectURL(ctx, url.ShortCode)
	assert.NoError(t, err)

	assert.Equal(t, url.OriginalURL, originalURL)
	assert.NotEmpty(t, originalURL)
	mockRepo.AssertExpectations(t)
}

func TestRedirectURL_EmptyShortCode(t *testing.T) {
	mockRepo, service := setupMockRepository()
	ctx := context.Background()

	mockRepo.On("Begin", ctx).Return(nil)
	mockRepo.On("FindByShortCode", ctx, "").Return(nil, errors.New("not found"))
	mockRepo.On("Rollback").Return(nil)

	originalURL, err := service.RedirectURL(context.Background(), "")
	assert.Error(t, err)

	assert.Empty(t, originalURL)
	assert.Contains(t, err.Error(), "failed to find URL: not found")
	mockRepo.AssertExpectations(t)
}

func TestRedirectURL_UpdateClickFail(t *testing.T) {
	mockRepo, service := setupMockRepository()
	ctx := context.Background()

	url := &models.URL{
		ID:        1,
		ShortCode: "abc123",
	}

	mockRepo.On("Begin", ctx).Return(nil)
	mockRepo.On("FindByShortCode", ctx, url.ShortCode).Return(url, nil)
	mockRepo.On("UpdateClickByID", ctx, url.ID).Return(errors.New("failed to update click count"))
	mockRepo.On("Rollback").Return(nil)

	originalURL, err := service.RedirectURL(ctx, url.ShortCode)
	assert.Error(t, err)

	assert.Empty(t, originalURL)
	assert.Contains(t, err.Error(), "failed to update click count")
	mockRepo.AssertExpectations(t)
}
