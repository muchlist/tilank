package crypt

import (
	"github.com/stretchr/testify/mock"
	"tilank/utils/rest_err"
)

type MockBcrypt struct {
	mock.Mock
}

func (m *MockBcrypt) GenerateHash(password string) (string, resterr.APIError) {
	args := m.Called(password)
	var err resterr.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(resterr.APIError)
	}

	return args.Get(0).(string), err
}

func (m *MockBcrypt) IsPWAndHashPWMatch(password string, hashPass string) bool {
	args := m.Called(password, hashPass)
	return args.Bool(0)
}
