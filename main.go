package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mfmayer/idabot/internal/llm"
)

var openaiApiKey string
var authorizedChatPartnerID string

func main() {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	openaiApiKey = os.Getenv("OPENAI_API_KEY")
	authorizedChatPartnerID = os.Getenv("AUTHORIZED_CHAT_PARTNER_ID")

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while creating discord session:", err)
		os.Exit(-1)
		return
	}

	llm := llm.NewLLM(dg.State.User.ID, authorizedChatPartnerID, openaiApiKey)

	dg.AddHandler(llm.)
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageDelete)

	err = dg.Open()
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Verbindung:", err)
		return
	}

	fmt.Println("Bot läuft. Drücke CTRL-C zum Beenden.")
	<-make(chan struct{})
}

func isDirectMessage(s *discordgo.Session, channelID string) bool {
	channel, err := s.Channel(channelID)
	if err != nil {
		fmt.Printf("Fehler beim Abrufen des Channels: %v\n", err)
		return false
	}

	return channel.Type == discordgo.ChannelTypeDM
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if m.Author.ID == s.State.User.ID {
	// 	return
	// }

	// if m.Author.ID != authorizedChatPartner || !isDirectMessage(s, m.ChannelID) {
	// 	return
	// }
	if m.Author.ID != authorizedChatPartner {
		return
	}

	fmt.Printf("create: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)

	if strings.HasPrefix(m.Content, "!gpt4 ") {
		// query := strings.TrimPrefix(m.Content, "!gpt4 ")

		// prompt := "User: " + query + "\nAssistant:"

		// completion, err := generateResponse(prompt, apiKey)
		// if err != nil {
		// 	s.ChannelMessageSend(m.ChannelID, "Es gab einen Fehler bei der Kommunikation mit der GPT-4 API.")
		// 	return
		// }

		// response := "GPT-4: " + completion
		// s.ChannelMessageSend(m.ChannelID, response)
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

