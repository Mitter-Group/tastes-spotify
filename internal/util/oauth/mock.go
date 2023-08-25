package oauth

import (
	"github.com/stretchr/testify/mock"
)

type OAuthMock struct {
	mock.Mock
}

func (mock *OAuthMock) GetAccessToken() (string, error) {
	args := mock.Called()
	return args.Get(0).(string), args.Error(1)
}
