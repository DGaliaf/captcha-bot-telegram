package main

import (
	"captcha-bot/app/internal/bot"
	"captcha-bot/app/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var tgToken = "TGTOKEN"

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Cannot read config file. Error: %v", err)
	}

	token, err := config.GetToken(tgToken)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Telegram Bot Token [%v] successfully obtained from env variable $TGTOKEN\n", token)

	tgbot := bot.NewBot(cfg, token)

	tgbot.StartHandlers()
	tgbot.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}
