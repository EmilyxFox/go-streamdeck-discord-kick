package main

import (
	"log"
	"time"

	"github.com/emilyxfox/go-streamdeck-discord-kick/cmd/actions"
	"github.com/emilyxfox/go-streamdeck-discord-kick/cmd/discord"
	"github.com/emilyxfox/go-streamdeck-sdk/streamdeck"
)

func main() {
	action := &actions.DisconncetAction{
		ActionConfig: streamdeck.ActionConfig{UUID: "com.emilyxfox.vckicker.disconnect"},
	}
	streamdeck.RegisterAction(action)

	go streamdeck.StartPlugin()

	var globalSettings map[string]any
	var err error
	for {
		globalSettings, err = streamdeck.GetGlobalSettings()
		if err == nil && globalSettings != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if token, ok := globalSettings["discordToken"].(string); ok && token != "" {
		if err := discord.DiscordLogin(token); err != nil {
			log.Printf("Error logging into Discord: %v", err)
		} else {
			log.Println("Discord bot logged in successfully at startup.")
		}
	} else {
		log.Println("No bot token found in global settings, waiting for manual login...")
	}

	select {}
}
