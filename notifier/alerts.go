package notifier

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(event models.Event, snapshotURL string, snapshot io.Reader, eventid string) {
	// Create copy of snapshot for each alerting method
	var snap []byte
	if snapshot != nil {
		snap, _ = io.ReadAll(snapshot)
	}
	if config.ConfigData.Alerts.Discord.Enabled {
		SendDiscordMessage(event, bytes.NewReader(snap), eventid)
	}
	if config.ConfigData.Alerts.Gotify.Enabled {
		SendGotifyPush(event, snapshotURL, eventid)
	}
	if config.ConfigData.Alerts.SMTP.Enabled {
		SendSMTP(event, bytes.NewReader(snap), eventid)
	}
	if config.ConfigData.Alerts.Telegram.Enabled {
		SendTelegramMessage(event, bytes.NewReader(snap), eventid)
	}
	if config.ConfigData.Alerts.Pushover.Enabled {
		SendPushoverMessage(event, bytes.NewReader(snap), eventid)
	}
	if config.ConfigData.Alerts.Nfty.Enabled {
		SendNftyPush(event, bytes.NewReader(snap), eventid)
	}
}

// Build notification based on template
func renderMessage(sourceTemplate string, event models.Event) string {
	// Assign Frigate URL to extra event fields
	event.Extra.LocalURL = config.ConfigData.Frigate.Server
	event.Extra.PublicURL = config.ConfigData.Frigate.PublicURL

	// MQTT uses CurrentZones, Web API uses Zones
	// Combine into one object to use regardless of connection method
	event.Zones = append(event.Zones, event.CurrentZones...)
	// Join zones into plain comma-separated string
	event.Extra.ZoneList = strings.Join(event.Zones, ", ")

	// If certain time format is provided, re-format date / time string
	eventTime := time.Unix(int64(event.StartTime), 0)
	event.Extra.FormattedTime = eventTime.String()
	if config.ConfigData.Alerts.General.TimeFormat != "" {
		event.Extra.FormattedTime = eventTime.Format(config.ConfigData.Alerts.General.TimeFormat)
	}

	// For Web API query, top-level top_score value is no longer used
	// So need to replace it with data.top_score value
	if event.TopScore == 0 {
		event.TopScore = event.Data.TopScore
	}
	// Calc TopScore percentage
	event.Extra.TopScorePercent = fmt.Sprintf("%v%%", int((event.TopScore * 100)))

	// Render template
	var tmpl *template.Template
	if sourceTemplate == "markdown" || sourceTemplate == "plaintext" || sourceTemplate == "html" {
		var templateFile = "./templates/" + sourceTemplate + ".template"
		tmpl = template.Must(template.ParseFiles(templateFile))
	}

	var renderedTemplate bytes.Buffer
	err := tmpl.Execute(&renderedTemplate, event)
	if err != nil {
		panic(err)
	}

	return renderedTemplate.String()

}
