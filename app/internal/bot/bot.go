package bot

import (
	"captcha-bot/app/internal/config"
	"captcha-bot/app/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"sync"
	"time"
)

type Bot struct {
	cfg         config.Config
	bot         *tb.Bot
	passedUsers sync.Map
}

func NewBot(cfg config.Config) Bot {
	var httpClient *http.Client
	if cfg.UseSocks5Proxy == "yes" {
		var err error
		httpClient, err = utils.InitSocks5Client(cfg)
		if err != nil {
			log.Fatalln(err)
		}
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	})
	if err != nil {
		log.Fatalf("Cannot start bot. Error: %v\n", err)
	}

	return Bot{
		cfg:         cfg,
		bot:         bot,
		passedUsers: sync.Map{},
	}
}

func (b Bot) StartHandlers() {
	b.bot.Handle(tb.OnUserJoined, b.challengeUser)
	b.bot.Handle(tb.OnCallback, b.passChallenge)

	b.bot.Handle("/healthz", func(m *tb.Message) {
		msg := "I'm OK"
		if _, err := b.bot.Send(m.Chat, msg); err != nil {
			log.Println(err)
		}
		log.Printf("Healthz request from user: %v\n in chat: %v", m.Sender, m.Chat)
	})
}

func (b Bot) Start() {
	log.Println("Bot started!")
	go func() {
		b.bot.Start()
	}()
}
