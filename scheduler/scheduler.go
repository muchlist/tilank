package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"tilank/service"
	"tilank/utils/logger"
	"time"
)

func RunScheduler(
	truckService *service.TruckService,
) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		logger.Error("gagal menggunakan timezone wita", err)
	}
	s := gocron.NewScheduler(witaTimeZone)

	// run blokir reset
	_, _ = s.Every(1).Hours().Do(func() {
		truckAffected, err := truckService.ResetBlockedTruck()
		if err != nil {
			logger.Error("Reset blokir truck error", err)
		}
		if truckAffected != 0 {
			logger.Info(fmt.Sprintf("Reset blokir truck diterapkan ke %d truck", truckAffected))
		}
	})

	s.StartAsync()
}
