package paydextoml

import "github.com/stretchr/testify/mock"

// MockClient is a mockable paydextoml client.
type MockClient struct {
	mock.Mock
}

// GetPaydexToml is a mocking a method
func (m *MockClient) GetPaydexToml(domain string) (*Response, error) {
	a := m.Called(domain)
	return a.Get(0).(*Response), a.Error(1)
}

// GetPaydexTomlByAddress is a mocking a method
func (m *MockClient) GetPaydexTomlByAddress(address string) (*Response, error) {
	a := m.Called(address)
	return a.Get(0).(*Response), a.Error(1)
}
