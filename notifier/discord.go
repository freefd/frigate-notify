package notifier

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

// SendDiscordMessage pushes alert message to Discord via webhook
func SendDiscordMessage(event models.Event, snapshot io.Reader, eventid string) {
	var err error

	// Build notification
	message := renderMessage("markdown", event)

	// Connect to Discord
	client, err := webhook.NewWithURL(config.ConfigData.Alerts.Discord.Webhook)
	if err != nil {
		log.Printf("Event ID %v - Unable to send Discord Alert: %v", eventid, err)
	}
	defer client.Close(context.TODO())

	title := fmt.Sprintf("**%v**\n\n", config.ConfigData.Alerts.General.Title)
	message = title + message

	// Send alert & attach snapshot if one was saved
	if snapshot != nil {
		image := discord.NewFile("snapshot.jpg", "", snapshot)
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetImage("attachment://snapshot.jpg").SetColor(5793266).Build()
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).SetFiles(image).Build())
	} else {
		message += "\nNo snapshot saved."
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetColor(5793266).Build()
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).Build())

	}
	if err != nil {
		log.Printf("Event ID %v - Unable to send Discord Alert: %v", eventid, err)
	}
	log.Printf("Event ID %v - Discord alert sent", eventid)
}
