package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

func (c *Chat) sender() {
	for msg := range c.messages {
		if _, err := c.connection.Write(c.formatMessage(msg)); err != nil {
			slog.Error("writer.write", slog.Any("error", err))
			return
		}
	}
}

func (c *Chat) formatMessage(message string) []byte {
	return []byte(fmt.Sprintf("%s: %s\n", c.name, message))
}

func (c *Chat) reader() {
	scanner := bufio.NewScanner(c.connection)
	for scanner.Scan() {
		c.history = append(c.history, scanner.Text())
	}
	slog.Info("connection is closed")
	os.Exit(0)
}
