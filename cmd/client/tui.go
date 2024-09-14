package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Chat) Init() tea.Cmd {
	return cursor.Blink
}

func (c *Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd
	c.textInput, tiCmd = c.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			c.byeMessage()
			return c, tea.Quit
		case tea.KeyEnter:
			message := c.textInput.Value()
			if message == "" {
				break
			}
			c.messages <- message
			c.textInput.Reset()
		}
	}
	return c, tiCmd
}

func (c *Chat) View() string {
	if len(c.history) == 0 {
		return fmt.Sprintf("%s\nTo leave the chat, press ESC or Ctrl+C", c.textInput.View())
	}

	return fmt.Sprintf(
		"%s\n%s\nTo leave the chat, press ESC or Ctrl+C",
		strings.Join(c.history, "\n"),
		c.textInput.View(),
	)
}
