package actions

import (
	"log"

	"github.com/emilyxfox/go-streamdeck-discord-kick/cmd/discord"
	"github.com/emilyxfox/go-streamdeck-sdk/streamdeck"
)

type DisconncetAction struct {
	streamdeck.ActionConfig
}

func (a *DisconncetAction) HandleKeyDown(event *streamdeck.KeyDownEvent) {
	settings, err := event.GetSettings()
	if err != nil {
		event.ShowAlert()
		log.Println("Error getting action settings:", err)
		return
	}
	userID, uOk := settings["user_id"].(string)
	guildID, gOk := settings["guild_id"].(string)
	if !uOk || !gOk {
		event.ShowAlert()
		log.Println("Action settings must include a valid user_id and guild_id")
		return
	}

	if err := discord.DisconnectFromVC(guildID, userID); err != nil {
		log.Printf("Failed to disconnect user %s: %v", userID, err)
		event.ShowAlert()
		return
	}

	event.ShowOk()
}

func (a *DisconncetAction) HandleSendToPlugin(event *streamdeck.SendToPluginEvent) {
	log.Printf("%v", event.Payload["request"])
	if event.Payload["request"] == "discordLogin" {
		log.Println("Trying to log into discord")
		globalSettings, err := event.GetGlobalSettings()
		if err != nil {
			log.Printf("Failed to get global settings")
			return
		}
		discordToken, ok := globalSettings["discordToken"].(string)
		if (discordToken) == "" || !ok {
			log.Println("Discord token is empty")
			return
		}

		if err := discord.DiscordLogin(discordToken); err != nil {
			log.Printf("Error logging into Discord: %v", err)
			return
		}
	}
}

func (a *DisconncetAction) HandleDidReceiveSettings(event *streamdeck.DidReceiveSettingsEvent) {
	log.Println("handling settings")
	settings := event.Payload.Settings
	userID, ok := settings["user_id"].(string)
	if !ok {
		log.Printf("Failed to get userID from settings")
		return
	}

	img, err := discord.GetUserAvatar(userID)
	if err != nil {
		log.Printf("Failed to get avatar for %v: %v", userID, err)
		return
	}

	if err := event.SetImage(img); err != nil {
		log.Println("Error setting action image:", err)
	}
}
