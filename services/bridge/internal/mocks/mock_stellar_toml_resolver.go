package mocks

import (
	"github.com/paydex-core/paydex-go/clients/paydextoml"
	"github.com/stretchr/testify/mock"
)

// MockPaydextomlResolver ...
type MockPaydextomlResolver struct {
	mock.Mock
}

// GetPaydexToml is a mocking a method
func (m *MockPaydextomlResolver) GetPaydexToml(domain string) (resp *paydextoml.Response, err error) {
	a := m.Called(domain)
	return a.Get(0).(*paydextoml.Response), a.Error(1)
}

// GetPaydexTomlByAddress is a mocking a method
func (m *MockPaydextomlResolver) GetPaydexTomlByAddress(addy string) (*paydextoml.Response, error) {
	a := m.Called(addy)
	return a.Get(0).(*paydextoml.Response), a.Error(1)
}
