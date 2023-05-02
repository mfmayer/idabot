package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/mfmayer/idabot/internal/dbot"
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

	err = dg.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while opening connection:", err)
		os.Exit(-1)
		return
	}

	bot := dbot.NewDBOT(dg.State.User.ID, authorizedChatPartnerID, openaiApiKey)
	dg.AddHandler(bot.DiscordMessageCreate)
	dg.AddHandler(bot.DiscordMessageDelete)
	dg.AddHandler(bot.DiscordMessageUpdate)

	fmt.Println("Bot running. Press CTRL-C to close.")
	<-make(chan struct{})
}
