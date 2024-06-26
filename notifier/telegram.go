package notifier

import (
	"io"
	"log"
	"strings"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendTelegramMessage sends alert through Telegram to individual users
func SendTelegramMessage(event models.Event, snapshot io.Reader, eventid string) {
	// Build notification
	message := renderMessage("html", event)
	message = strings.ReplaceAll(message, "<br />", "")

	bot, err := tgbotapi.NewBotAPI(config.ConfigData.Alerts.Telegram.Token)
	if err != nil {
		log.Printf("Event ID %v - Failed to connect to Telegram: %v", eventid, err)
		return
	}

	if snapshot != nil {
		// Attach & send snapshot
		photo := tgbotapi.NewPhoto(config.ConfigData.Alerts.Telegram.ChatID, tgbotapi.FileReader{Name: "Snapshot", Reader: snapshot})
		photo.Caption = message
		photo.ParseMode = "HTML"
		if _, err := bot.Send(photo); err != nil {
			log.Printf("Event ID %v - Failed to send alert via Telegram: %v", eventid, err)
			return
		}
	} else {
		// Send plain text message if no snapshot available
		message += "No snapshot saved."
		msg := tgbotapi.NewMessage(config.ConfigData.Alerts.Telegram.ChatID, message)
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Event ID %v - Failed to send alert via Telegram: %v", eventid, err)
			return
		}
	}
	log.Printf("Event ID %v - Telegram alert sent", eventid)
}
