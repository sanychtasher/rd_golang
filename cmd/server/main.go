package main

import (
	"flag"
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	slog.Info("server is running")

	var opt ServerOptions

	flag.StringVar(&opt.Host, "host", ":8456", "server host")
	flag.Parse()

	if err := NewServer(opt).Run(opt.Host); err != nil {
		slog.Error("new_server", slog.Any("error", err))
		return
	}

	select {}
}
