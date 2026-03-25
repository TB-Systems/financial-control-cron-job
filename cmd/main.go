package main

import (
	"auto-service/internal/services"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	service := services.NewRecurringTransactionsService()
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger), // evita overlap do mesmo job
			cron.Recover(cron.DefaultLogger),
		),
	)

	_, _ = c.AddFunc("0 0 2 * * *", func() {
		if err := service.Run(); err != nil {
			log.Printf("job erro: %v", err)
		}
	})

	c.Start()
	select {}
}
