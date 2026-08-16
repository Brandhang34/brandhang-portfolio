package handler

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	maxSubjectLen = 200
	maxMessageLen = 4000
)

// ErrInvalidContactInput is returned when a contact submission fails validation.
var ErrInvalidContactInput = errors.New("invalid contact submission")

// ContactMsg is a validated contact-form submission.
type ContactMsg struct {
	Email   string
	Subject string
	Message string
}

// NewContactMsg validates raw form input and caps field lengths so the Discord
// relay can't be used to post arbitrarily large payloads.
func NewContactMsg(email, subject, message string) (ContactMsg, error) {
	email = strings.TrimSpace(email)
	subject = strings.TrimSpace(subject)
	message = strings.TrimSpace(message)

	if _, err := mail.ParseAddress(email); err != nil {
		return ContactMsg{}, fmt.Errorf("%w: email", ErrInvalidContactInput)
	}
	if message == "" {
		return ContactMsg{}, fmt.Errorf("%w: empty message", ErrInvalidContactInput)
	}
	if len(subject) > maxSubjectLen {
		return ContactMsg{}, fmt.Errorf("%w: subject too long", ErrInvalidContactInput)
	}
	if len(message) > maxMessageLen {
		return ContactMsg{}, fmt.Errorf("%w: message too long", ErrInvalidContactInput)
	}

	return ContactMsg{Email: email, Subject: subject, Message: message}, nil
}

// SendDiscordMsg posts a contact submission to the configured Discord channel.
//
// ChannelMessageSend is a plain REST call, so there is deliberately no
// discord.Open() here -- opening a gateway websocket per submission is both
// expensive and a fast route to a rate-limited bot.
func SendDiscordMsg(msg ContactMsg) error {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if token == "" || channelID == "" {
		return errors.New("discord: DISCORD_BOT_TOKEN and DISCORD_CHANNEL_ID must be set")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("discord: new session: %w", err)
	}

	body := fmt.Sprintf("**From:** %s\n**Subject:** %s\n**Message:**\n%s",
		msg.Email, msg.Subject, msg.Message)

	if _, err := discord.ChannelMessageSend(channelID, body); err != nil {
		return fmt.Errorf("discord: send message: %w", err)
	}

	return nil
}
