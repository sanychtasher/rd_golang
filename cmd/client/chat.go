package main

import (
	"fmt"
	"net"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Chat{}

type Chat struct {
	textInput  textinput.Model
	name       string
	host       string
	history    []string
	messages   chan string
	connection net.Conn
}

type ChatOptions struct {
	Name string
	Host string
}

func NewChat(opt *ChatOptions) (*Chat, error) {
	textInput := textinput.New()
	textInput.Focus()
	textInput.CharLimit = 20

	c := &Chat{
		textInput: textInput,
		name:      opt.Name,
		host:      opt.Host,
		history:   make([]string, 0),
		messages:  make(chan string),
	}

	var err error
	c.connection, err = net.Dial("tcp", c.host)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	go c.reader()
	go c.sender()

	c.helloMessage()

	return c, nil
}

func (c *Chat) helloMessage() {
	c.messages <- "Hello Everyone!"
}

func (c *Chat) byeMessage() {
	c.messages <- "Bye Everyone!"
}
