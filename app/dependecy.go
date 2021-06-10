package app

import (
	"tilank/dao/user_dao"
	"tilank/handler"
	"tilank/service"
	"tilank/utils/crypt"
	"tilank/utils/mjwt"
)

var (
	// Utils
	cryptoUtils = crypt.NewCrypto()
	jwt         = mjwt.NewJwt()

	// Dao
	userDao = userdao.NewUserDao()

	// Service
	userService = service.NewUserService(userDao, cryptoUtils, jwt)

	// Controller or Handler
	pingHandler = handler.NewPingHandler()
	userHandler = handler.NewUserHandler(userService)
)
