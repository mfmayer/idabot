package dbot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/mfmayer/goai"
)

type DBot struct {
	id            string
	partnerID     string
	oai           *goai.Client
	conversations map[string]*Conversation
}

// func isDirectMessage(s *discordgo.Session, channelID string) bool {
// 	channel, err := s.Channel(channelID)
// 	if err != nil {
// 		fmt.Printf("Fehler beim Abrufen des Channels: %v\n", err)
// 		return false
// 	}
// 	return channel.Type == discordgo.ChannelTypeDM
// }

// NewLLM creates new llm bot
func NewDBOT(id string, partnerID string, openaiApiKey string) *DBot {
	return &DBot{
		id:            id,
		partnerID:     partnerID,
		conversations: map[string]*Conversation{},
		oai:           goai.NewClient(openaiApiKey),
	}
}

func (bot *DBot) Conversation(channelID string) *Conversation {
	c, ok := bot.conversations[channelID]
	if !ok {
		c = NewConversation()
		bot.conversations[channelID] = c
	}
	return c
}

func (bot *DBot) getContext(authorID string, channelID string) (followUp bool, role goai.Role, conversation *Conversation) {
	followUp = true
	switch authorID {
	case bot.partnerID:
		role = goai.RoleUser
	case bot.id:
		role = goai.RoleAssistant
	default:
		followUp = false
		return
	}
	conversation = bot.Conversation(channelID)
	return
}

func (bot *DBot) DiscordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	followUp, role, conversation := bot.getContext(m.Author.ID, m.ChannelID)
	if !followUp {
		return
	}
	fmt.Printf("create: from:%s (%v-%v): %v\n", role, m.ChannelID, m.ID, m.Content)
	conversation.AddMessage(m.ID, role, m.Content)
	// In case user wrote the message, create chat prompt and get completion from open ai client
	if role == "user" {
		prompt := conversation.NewChatPrompt()
		fmt.Println(prompt)
		if chatCompletion, err := bot.oai.GetChatCompletion(prompt); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			if len(chatCompletion.Choices) > 0 {
				s.ChannelMessageSend(m.ChannelID, chatCompletion.Choices[0].Message.Content)
			}
		}
	}
}

func (bot *DBot) DiscordMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("delete: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)
}

func (bot *DBot) DiscordMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	fmt.Printf("update: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)
}
