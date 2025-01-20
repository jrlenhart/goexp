package limiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int, error) {
	args := m.Called(ctx, key)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) Expire(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockStorage) BlockKey(ctx context.Context, key string, duration time.Duration) error {
	args := m.Called(ctx, key, duration)
	return args.Error(0)
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func TestRateLimiter_AllowRequest_Increment_Expire(t *testing.T) {
	storage := new(MockStorage)
	rateLimiter := NewRateLimiter(storage, 3, 5, 10, 20)

	storage.On("IsBlocked", mock.Anything, "ip:192.168.0.1").Return(false, nil)
	storage.On("Increment", mock.Anything, "ip:192.168.0.1").Return(1, nil)
	storage.On("Expire", mock.Anything, "ip:192.168.0.1", time.Second).Return(nil)

	err := rateLimiter.AllowRequest(context.Background(), "ip:192.168.0.1")

	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestRateLimiter_AllowRequest_Increment(t *testing.T) {
	storage := new(MockStorage)
	rateLimiter := NewRateLimiter(storage, 3, 5, 10, 20)

	storage.On("IsBlocked", mock.Anything, "ip:192.168.0.1").Return(false, nil)
	storage.On("Increment", mock.Anything, "ip:192.168.0.1").Return(2, nil)

	err := rateLimiter.AllowRequest(context.Background(), "ip:192.168.0.1")

	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestRateLimiter_AllowRequest_ExceedLimitAndBlock(t *testing.T) {
	storage := new(MockStorage)
	rateLimiter := NewRateLimiter(storage, 3, 5, 10, 20)

	storage.On("IsBlocked", mock.Anything, "ip:192.168.0.1").Return(false, nil)
	storage.On("Increment", mock.Anything, "ip:192.168.0.1").Return(rateLimiter.maxRequestIP+1, nil)
	storage.On("BlockKey", mock.Anything, "ip:192.168.0.1", time.Duration(rateLimiter.blockTimeIP)*time.Second).Return(nil)

	err := rateLimiter.AllowRequest(context.Background(), "ip:192.168.0.1")

	assert.Error(t, err)
	assert.Equal(t, "key ip:192.168.0.1 exceeded the limit and is now blocked", err.Error())
	storage.AssertExpectations(t)

	storage.On("IsBlocked", mock.Anything, "token:abc123").Return(false, nil)
	storage.On("Increment", mock.Anything, "token:abc123").Return(rateLimiter.maxRequestToken+1, nil)
	storage.On("BlockKey", mock.Anything, "token:abc123", time.Duration(rateLimiter.blockTimeToken)*time.Second).Return(nil)

	err = rateLimiter.AllowRequest(context.Background(), "token:abc123")

	assert.Error(t, err)
	assert.Equal(t, "key token:abc123 exceeded the limit and is now blocked", err.Error())
	storage.AssertExpectations(t)
}

func TestRateLimiter_AllowRequest_KeyBlocked(t *testing.T) {
	storage := new(MockStorage)
	rateLimiter := NewRateLimiter(storage, 3, 5, 10, 20)

	storage.On("IsBlocked", mock.Anything, "ip:192.168.0.1").Return(true, nil)

	err := rateLimiter.AllowRequest(context.Background(), "ip:192.168.0.1")

	assert.Error(t, err)
	assert.Equal(t, "key ip:192.168.0.1 is blocked", err.Error())
	storage.AssertExpectations(t)
}
