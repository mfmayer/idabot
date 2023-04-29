package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type completionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type LLM struct {
	id           string
	partnerID    string
	openaiApiKey string
}

// NewLLM creates new llm bot
func NewLLM(id string, partnerID string, openaiApiKey string) *LLM {
	return &LLM{
		id:           id,
		openaiApiKey: openaiApiKey,
		// conversations: map[string]
	}
}

func (lm *LLM) DiscordMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("delete: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)
}

func (lm *LLM) DiscordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("delete: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)
	if m.Author.ID != lm.partnerID {
		return
	}
	if strings.HasPrefix(m.Content, "!gpt4") {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func (lm *LLM) DiscordMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	fmt.Printf("update: %v-%v: %v\n", m.ChannelID, m.ID, m.Content)
}

func (lm *LLM) generateResponse(prompt, apiKey string) (string, error) {
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

	var completionResponse completionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", err
	}

	return completionResponse.Choices[0].Text, nil
}
