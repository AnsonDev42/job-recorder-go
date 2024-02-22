package utils

import (
	"github.com/go-co-op/gocron"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"log"
	"strconv"
)

func SetupTelegramBot(apiToken string, receiverID string) (err error) {
	rID, err := strconv.ParseInt(receiverID, 10, 64)
	if err != nil {
		return err
	}
	telegramService, _ := telegram.New(apiToken)
	telegramService.AddReceivers(rID)
	notify.UseServices(telegramService)

	return nil
}

func SetSummaryScheduler(s *gocron.Scheduler, job func()) {
	// Send a summary every day at 20:00
	log.Print("Setting up the daily summary scheduler for 20:00")
	_, err := s.Every(1).Day().At("20:00").Do(job)
	if err != nil {
		panic("failed to schedule the daily summary")
	}
}
