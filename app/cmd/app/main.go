package main

import (
	"captcha-bot/app/internal/bot"
	"captcha-bot/app/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Cannot read config file. Error: %v", err)
	}

	tgbot := bot.NewBot(cfg)

	tgbot.Start()
	tgbot.StartHandlers()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}
