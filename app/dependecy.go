package app

import (
	"tilank/dao/jptdao"
	"tilank/dao/userdao"
	"tilank/dao/violationdao"
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
	userDao      = userdao.NewUserDao()
	violationDao = violationdao.NewViolationDao()
	jptDao       = jptdao.NewJptDao()

	// Service
	userService      = service.NewUserService(userDao, cryptoUtils, jwt)
	jptService       = service.NewJptService(jptDao)
	violationService = service.NewViolationService(violationDao, jptDao)

	// Controller or Handler
	pingHandler      = handler.NewPingHandler()
	userHandler      = handler.NewUserHandler(userService)
	violationHandler = handler.NewViolationHandler(violationService)
)
