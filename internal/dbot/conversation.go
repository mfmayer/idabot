package dbot

import (
	"github.com/mfmayer/goai"
)

type Message struct {
	oaiMessage  *goai.Message
	nextMessage *Message
}

func NewMessage(role goai.Role, content string) *Message {
	m := &Message{}
	m.Update(role, content)
	return m
}

func (m *Message) Update(role goai.Role, content string) {
	if m.oaiMessage == nil {
		m.oaiMessage = &goai.Message{}
	}
	m.oaiMessage.Content = content
	m.oaiMessage.Role = role
}

func (m *Message) SetNext(next *Message) {
	m.nextMessage = next
}

func (m *Message) Next() *Message {
	return m.nextMessage
}

type Conversation struct {
	first    *Message
	last     *Message
	messages map[string]*Message
}

func NewConversation() *Conversation {
	return &Conversation{
		messages: map[string]*Message{},
	}
}

func (c *Conversation) AddMessage(id string, role goai.Role, content string) {
	m := NewMessage(role, content)
	if c.first == nil {
		c.first = m
	}
	if c.last != nil {
		c.last.SetNext(m)
		c.last = m
	}
	c.last = m
}

func (c *Conversation) UpdateMessage(id string, role goai.Role, content string) {
	m := c.messages[id]
	if m == nil {
		return
	}
	m.Update(role, content)
}

func (c *Conversation) NewChatPrompt() *goai.ChatPrompt {
	chatPrompt := &goai.ChatPrompt{}
	for msg := c.first; msg != nil; msg = msg.Next() {
		chatPrompt.Messages = append(chatPrompt.Messages, msg.oaiMessage)
	}
	return chatPrompt
}
