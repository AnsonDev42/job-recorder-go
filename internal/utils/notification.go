package utils

import (
	"context"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
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

	// Send a test message.
	_ = notify.Send(
		context.Background(),
		"Subject/Title",
		"The actual message - Hello, you awesome gophers! :)",
	)
	return nil
}
