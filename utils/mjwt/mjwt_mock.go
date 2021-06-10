package mjwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
	"tilank/utils/rest_err"
)

type MockJwt struct {
	mock.Mock
}

func (m *MockJwt) GenerateToken(claims CustomClaim) (string, resterr.APIError) {
	args := m.Called(claims)
	var err resterr.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(resterr.APIError)
	}

	return args.Get(0).(string), err
}

func (m *MockJwt) ValidateToken(tokenString string) (*jwt.Token, resterr.APIError) {
	args := m.Called(tokenString)
	var res *jwt.Token
	if args.Get(0) != nil {
		res = args.Get(0).(*jwt.Token)
	}

	var err resterr.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(resterr.APIError)
	}

	return res, err
}

func (m *MockJwt) ReadToken(token *jwt.Token) (*CustomClaim, resterr.APIError) {
	args := m.Called(token)
	var res *CustomClaim
	if args.Get(0) != nil {
		res = args.Get(0).(*CustomClaim)
	}

	var err resterr.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(resterr.APIError)
	}

	return res, err
}
