package main

import (
	"flag"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	opt := new(ChatOptions)

	flag.StringVar(&opt.Name, "name", "default", "username")
	flag.StringVar(&opt.Host, "host", "localhost:8456", "server host")
	flag.Parse()

	chat, err := NewChat(opt)
	if err != nil {
		slog.Error("new_chat", slog.Any("error", err))
		return
	}

	if _, err = tea.NewProgram(chat).Run(); err != nil {
		slog.Error("run", slog.Any("error", err))
		return
	}
}
