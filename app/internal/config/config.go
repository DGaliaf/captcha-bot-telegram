package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	BotToken            string `toml:"bot_token"`
	ButtonText          string `toml:"button_text"`
	WelcomeMessage      string `toml:"welcome_message"`
	AfterSuccessMessage string `toml:"after_success_message"`
	AfterFailMessage    string `toml:"after_fail_message"`
	PrintSuccessAndFail string `toml:"print_success_and_fail_messages_strategy"`
	WelcomeTimeout      string `toml:"welcome_timeout"`
	BanDurations        string `toml:"ban_duration"`
	UseSocks5Proxy      string `toml:"use_socks5_proxy"`
	Socks5Address       string `toml:"socks5_address"`
	Socks5Port          string `toml:"socks5_port"`
	Socks5Login         string `toml:"socks5_login"`
	Socks5Password      string `toml:"socks5_password"`
}

const CONFIGPATH = "config.toml"

var config Config

func NewConfig() (Config, error) {
	err := readConfig()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func readConfig() (err error) {
	_, err = toml.DecodeFile(CONFIGPATH, &config)
	if err != nil {
		return err
	}

	return
}
