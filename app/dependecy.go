package app

import (
	"tilank/clients/fcm"
	"tilank/dao/jptdao"
	"tilank/dao/rulesdao"
	"tilank/dao/truckdao"
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
	truckDao     = truckdao.NewTruckDao()
	rulesDao     = rulesdao.NewRulesDao()

	// api client
	fcmClient = fcm.NewFcmClient()

	// Service
	userService      = service.NewUserService(userDao, cryptoUtils, jwt)
	jptService       = service.NewJptService(jptDao)
	violationService = service.NewViolationService(violationDao, truckDao, rulesDao, userDao, fcmClient)
	truckService     = service.NewTruckService(truckDao)
	rulesService     = service.NewRulesService(rulesDao)

	// Controller or Handler
	pingHandler      = handler.NewPingHandler()
	userHandler      = handler.NewUserHandler(userService)
	violationHandler = handler.NewViolationHandler(violationService)
	jptHandler       = handler.NewJptHandler(jptService)
	truckHandler     = handler.NewTruckHandler(truckService)
	rulesHandler     = handler.NewRulesHandler(rulesService)
)
