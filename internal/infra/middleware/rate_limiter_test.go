package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/carlosmeds/rate-limiter/configs"
)

func TestGetCredentials(t *testing.T) {
	tests := []struct {
		name         string
		apiKey       string
		remoteAddr   string
		expectedKey  string
		expectedAddr string
	}{
		{
			name:         "Valid API Key and Remote Address",
			apiKey:       "test-api-key",
			remoteAddr:   "192.168.1.1",
			expectedKey:  "test-api-key",
			expectedAddr: "192.168.1.1",
		},
		{
			name:         "Empty API Key",
			apiKey:       "",
			remoteAddr:   "192.168.1.1",
			expectedKey:  "",
			expectedAddr: "192.168.1.1",
		},
		{
			name:         "Empty Remote Address",
			apiKey:       "test-api-key",
			remoteAddr:   "",
			expectedKey:  "test-api-key",
			expectedAddr: "",
		},
		{
			name:         "Empty API Key and Remote Address",
			apiKey:       "",
			remoteAddr:   "",
			expectedKey:  "",
			expectedAddr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("API_KEY", tt.apiKey)
			req.RemoteAddr = tt.remoteAddr

			apiKey, clientIP := getCredentials(req)
			if apiKey != tt.expectedKey {
				t.Errorf("Expected API Key: %v, got: %v", tt.expectedKey, apiKey)
			}
			if clientIP != tt.expectedAddr {
				t.Errorf("Expected Remote Address: %v, got: %v", tt.expectedAddr, clientIP)
			}
		})
	}
}

func TestGetRequestsKey(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		clientIP    string
		expectedKey string
	}{
		{
			name:        "Valid API Key and Client IP",
			apiKey:      "test-api-key",
			clientIP:    "192.168.1.1",
			expectedKey: "requests@test-api-key",
		},
		{
			name:        "Empty API Key",
			apiKey:      "",
			clientIP:    "192.168.1.1",
			expectedKey: "requests@192.168.1.1",
		},
		{
			name:        "Empty Client IP",
			apiKey:      "test-api-key",
			clientIP:    "",
			expectedKey: "requests@test-api-key",
		},
		{
			name:        "Empty API Key and Client IP",
			apiKey:      "",
			clientIP:    "",
			expectedKey: "requests@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRequestsKey(tt.apiKey, tt.clientIP)
			if result != tt.expectedKey {
				t.Errorf("Expected Key: %v, got: %v", tt.expectedKey, result)
			}
		})
	}
}
func TestGetBlackListKey(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		clientIP    string
		expectedKey string
	}{
		{
			name:        "Valid API Key and Client IP",
			apiKey:      "test-api-key",
			clientIP:    "192.168.1.1",
			expectedKey: "blacklist@test-api-key",
		},
		{
			name:        "Empty API Key",
			apiKey:      "",
			clientIP:    "192.168.1.1",
			expectedKey: "blacklist@192.168.1.1",
		},
		{
			name:        "Empty Client IP",
			apiKey:      "test-api-key",
			clientIP:    "",
			expectedKey: "blacklist@test-api-key",
		},
		{
			name:        "Empty API Key and Client IP",
			apiKey:      "",
			clientIP:    "",
			expectedKey: "blacklist@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBlackListKey(tt.apiKey, tt.clientIP)
			if result != tt.expectedKey {
				t.Errorf("Expected Key: %v, got: %v", tt.expectedKey, result)
			}
		})
	}
}

func TestIsBlackListed(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		blackListed  string
		getErr       error
		expectedMsg  string
		expectedCode int
	}{
		{
			name:         "Not Blacklisted",
			key:          "blacklist@test-api-key",
			blackListed:  "",
			getErr:       nil,
			expectedMsg:  "",
			expectedCode: 0,
		},
		{
			name:         "Blacklisted",
			key:          "blacklist@test-api-key",
			blackListed:  "Too many requests",
			getErr:       nil,
			expectedMsg:  rateLimitMsg,
			expectedCode: http.StatusTooManyRequests,
		},
		{
			name:         "Error in Get",
			key:          "blacklist@test-api-key",
			blackListed:  "",
			getErr:       fmt.Errorf("some error"),
			expectedMsg:  internalErrMsg,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockStore := &MockStore{
				GetFunc: func(ctx context.Context, key string) (string, error) {
					if key != tt.key {
						t.Errorf("Expected key: %v, got: %v", tt.key, key)
					}
					return tt.blackListed, tt.getErr
				},
			}
			md := &RateLimiterMiddleware{s: mockStore}

			msg, code := md.isBlackListed(ctx, tt.key)
			if msg != tt.expectedMsg {
				t.Errorf("Expected message: %v, got: %v", tt.expectedMsg, msg)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected code: %v, got: %v", tt.expectedCode, code)
			}
		})
	}
}

func TestGetLimit(t *testing.T) {
	tests := []struct {
		name          string
		apiKey        string
		defaultLimit  int64
		apiKeyLimits  map[string]int64
		expectedLimit int64
		expectedMsg   string
		expectedCode  int
	}{
		{
			name:          "Empty API Key",
			apiKey:        "",
			defaultLimit:  100,
			apiKeyLimits:  map[string]int64{},
			expectedLimit: 100,
			expectedMsg:   "",
			expectedCode:  0,
		},
		{
			name:          "Valid API Key",
			apiKey:        "test-api-key",
			defaultLimit:  100,
			apiKeyLimits:  map[string]int64{"test-api-key": 200},
			expectedLimit: 200,
			expectedMsg:   "",
			expectedCode:  0,
		},
		{
			name:          "Invalid API Key",
			apiKey:        "invalid-api-key",
			defaultLimit:  100,
			apiKeyLimits:  map[string]int64{"test-api-key": 200},
			expectedLimit: 0,
			expectedMsg:   invalidKey,
			expectedCode:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := &configs.Config{
				DefaultLimit: tt.defaultLimit,
				ApiKeyLimits: tt.apiKeyLimits,
			}

			mockStore := &MockStore{}
			md := &RateLimiterMiddleware{s: mockStore}

			limit, msg, code := md.getLimit(tt.apiKey, mockConfig)
			if limit != tt.expectedLimit {
				t.Errorf("Expected limit: %v, got: %v", tt.expectedLimit, limit)
			}
			if msg != tt.expectedMsg {
				t.Errorf("Expected message: %v, got: %v", tt.expectedMsg, msg)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected code: %v, got: %v", tt.expectedCode, code)
			}
		})
	}
}

func TestGetReachedLimit(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		limit         int64
		reachedLimit  bool
		hasReachedErr error
		expectedMsg   string
		expectedCode  int
	}{
		{
			name:          "Limit Not Reached",
			key:           "requests@test-api-key",
			limit:         100,
			reachedLimit:  false,
			hasReachedErr: nil,
			expectedMsg:   "",
			expectedCode:  0,
		},
		{
			name:          "Limit Reached",
			key:           "requests@test-api-key",
			limit:         100,
			reachedLimit:  true,
			hasReachedErr: nil,
			expectedMsg:   rateLimitMsg,
			expectedCode:  http.StatusTooManyRequests,
		},
		{
			name:          "Error in HasReachedLimit",
			key:           "requests@test-api-key",
			limit:         100,
			reachedLimit:  false,
			hasReachedErr: fmt.Errorf("some error"),
			expectedMsg:   internalErrMsg,
			expectedCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockStore := &MockStore{
				HasReachedLimitFunc: func(ctx context.Context, key string, limit int64) (bool, error) {
					if key != tt.key {
						t.Errorf("Expected key: %v, got: %v", tt.key, key)
					}
					if limit != tt.limit {
						t.Errorf("Expected limit: %v, got: %v", tt.limit, limit)
					}
					return tt.reachedLimit, tt.hasReachedErr
				},
			}
			md := &RateLimiterMiddleware{s: mockStore}

			msg, code := md.getReachedLimit(ctx, tt.key, tt.limit)
			if msg != tt.expectedMsg {
				t.Errorf("Expected message: %v, got: %v", tt.expectedMsg, msg)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected code: %v, got: %v", tt.expectedCode, code)
			}
		})
	}
}

func TestAddToBlackList(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		saveErr     error
		expectedErr error
	}{
		{
			name:        "Successfully added to blacklist",
			key:         "blacklist@test-api-key",
			saveErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "Error adding to blacklist",
			key:         "blacklist@test-api-key",
			saveErr:     fmt.Errorf("some error"),
			expectedErr: fmt.Errorf("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockConfig := &configs.Config{
				BlockedTime: 300,
			}
			mockStore := &MockStore{
				SaveFunc: func(ctx context.Context, key, value string, ttl int64) error {
					if key != tt.key {
						t.Errorf("Expected key: %v, got: %v", tt.key, key)
					}
					if value != "Too many requests" {
						t.Errorf("Expected value: %v, got: %v", "Too many requests", value)
					}
					return tt.saveErr
				},
			}
			md := &RateLimiterMiddleware{s: mockStore}

			err := md.AddToBlackList(ctx, tt.key, mockConfig)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if err == nil && tt.expectedErr != nil {
				t.Errorf("Expected error: %v, got no error", tt.expectedErr)
			}
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}
		})
	}
}

// MockStore is a mock implementation of the store interface used for testing
type MockStore struct {
	GetFunc             func(ctx context.Context, key string) (string, error)
	HasReachedLimitFunc func(ctx context.Context, key string, limit int64) (bool, error)
	SaveFunc            func(ctx context.Context, key, value string, ttl int64) error
}

func (m *MockStore) Get(ctx context.Context, key string) (string, error) {
	return m.GetFunc(ctx, key)
}

func (m *MockStore) HasReachedLimit(ctx context.Context, key string, limit int64) (bool, error) {
	return m.HasReachedLimitFunc(ctx, key, limit)
}

func (m *MockStore) Save(ctx context.Context, key, value string, ttl int64) error {
	return m.SaveFunc(ctx, key, value, ttl)
}
