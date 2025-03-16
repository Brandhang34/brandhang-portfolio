package handler

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func SendDiscordMsg(message string) {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	_, err = discord.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer discord.Close()

}
