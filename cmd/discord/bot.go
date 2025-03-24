package discord

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"

	"github.com/bwmarrin/discordgo"
)

var DiscordBot *discordgo.Session

func DiscordLogin(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("error logging into Discord bot: %v", err)
	}

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection to Discord: %v", err)
	}

	DiscordBot = dg
	log.Println("Discord bot logged in.")
	return nil
}

func DisconnectFromVC(guildID, userID string) error {
	if DiscordBot == nil {
		return fmt.Errorf("discord bot not logged in")
	}

	voiceState, err := DiscordBot.State.VoiceState(guildID, userID)
	log.Println(voiceState)
	if err != nil || voiceState.ChannelID == "" {
		return fmt.Errorf("user %s is not in a voice channel", userID)
	}

	err = DiscordBot.GuildMemberMove(guildID, userID, nil)
	if err != nil {
		return fmt.Errorf("error disconnecting user %s: %v", userID, err)
	}

	log.Printf("User %s disconnected from voice channel.", userID)
	return nil
}

func GetUserAvatar(userID string) (string, error) {
	if DiscordBot == nil {
		return "", fmt.Errorf("discord bot is not logged in")
	}
	user, err := DiscordBot.User(userID)
	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}
	img, err := DiscordBot.UserAvatarDecode(user)
	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("error encoding avatar image: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return fmt.Sprintf("data:image/png;base64,%s", base64Image), nil
}
