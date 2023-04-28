package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var openaiApiKey string
var authorizedChatPartner string

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	openaiApiKey = os.Getenv("OPENAI_API_KEY")
	authorizedChatPartner = os.Getenv("AUTHORIZED_CHAT_PARTNER")

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Fehler beim Erstellen des Discord-Session:", err)
		return
	}

	dg.AddHandler(messageCreate)

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

	if m.Author.ID != authorizedChatPartner || !isDirectMessage(s, m.ChannelID) {
		return
	}

	fmt.Printf("%v: %v\n", m.ChannelID, m.Content)

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

func generateResponse(prompt, apiKey string) (string, error) {
	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	data := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  150,
		"temperature": 0.8,
		"top_p":       1,
		"n":           1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", err
	}

	return completionResponse.Choices[0].Text, nil
}
