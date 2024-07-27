package bot

import (
	"captcha-bot/app/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

func (b Bot) challengeUser(m *tb.Message) {
	if m.UserJoined.ID != m.Sender.ID {
		return
	}
	log.Printf("User: %v joined the chat: %v", m.UserJoined, m.Chat)

	if member, err := b.bot.ChatMemberOf(m.Chat, m.UserJoined); err == nil {
		if member.Role == tb.Restricted {
			log.Printf("User: %v already restricted in chat: %v", m.UserJoined, m.Chat)
			return
		}
	}

	newChatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: false}}
	err := b.bot.Restrict(m.Chat, &newChatMember)
	if err != nil {
		log.Println(err)
	}

	inlineKeys := [][]tb.InlineButton{{tb.InlineButton{
		Unique: "challenge_btn",
		Text:   b.cfg.ButtonText,
	}}}

	challengeMsg, err := b.bot.Reply(m, b.cfg.WelcomeMessage, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	if err != nil {
		log.Printf("Can't send challenge msg: %v", err)
		return
	}

	n, err := strconv.ParseInt(b.cfg.WelcomeTimeout, 10, 64)
	if err != nil {
		log.Println(err)
	}
	time.AfterFunc(time.Duration(n)*time.Second, func() {
		_, passed := b.passedUsers.Load(m.UserJoined.ID)
		if !passed {
			banDuration, e := utils.GetBanDuration(b.cfg)
			if e != nil {
				log.Println(e)
			}
			chatMember := tb.ChatMember{User: m.UserJoined, RestrictedUntil: banDuration}
			err := b.bot.Ban(m.Chat, &chatMember)
			if err != nil {
				log.Println(err)
			}

			if b.cfg.PrintSuccessAndFail == "show" {
				_, err := b.bot.Edit(challengeMsg, b.cfg.AfterFailMessage)
				if err != nil {
					log.Println(err)
				}
			} else if b.cfg.PrintSuccessAndFail == "del" {
				err := b.bot.Delete(m)
				if err != nil {
					log.Println(err)
				}
				err = b.bot.Delete(challengeMsg)
				if err != nil {
					log.Println(err)
				}
			}

			log.Printf("User: %v was banned in chat: %v for: %v minutes", m.UserJoined, m.Chat, b.cfg.BanDurations)
		}
		b.passedUsers.Delete(m.UserJoined.ID)
	})
}

func (b Bot) passChallenge(c *tb.Callback) {
	if c.Message.ReplyTo.Sender.ID != c.Sender.ID {
		err := b.bot.Respond(c, &tb.CallbackResponse{Text: "This button isn't for you"})
		if err != nil {
			log.Println(err)
		}
		return
	}
	b.passedUsers.Store(c.Sender.ID, struct{}{})

	if b.cfg.PrintSuccessAndFail == "show" {
		_, err := b.bot.Edit(c.Message, b.cfg.AfterSuccessMessage)
		if err != nil {
			log.Println(err)
		}
	} else if b.cfg.PrintSuccessAndFail == "del" {
		err := b.bot.Delete(c.Message)
		if err != nil {
			log.Println(err)
		}
	}

	log.Printf("User: %v passed the challenge in chat: %v", c.Sender, c.Message.Chat)
	newChatMember := tb.ChatMember{User: c.Sender, RestrictedUntil: tb.Forever(), Rights: tb.Rights{CanSendMessages: true}}
	err := b.bot.Promote(c.Message.Chat, &newChatMember)
	if err != nil {
		log.Println(err)
	}
	err = b.bot.Respond(c, &tb.CallbackResponse{Text: "Validation passed!"})
	if err != nil {
		log.Println(err)
	}
}
