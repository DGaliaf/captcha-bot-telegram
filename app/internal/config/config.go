package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

type Config struct {
	BotToken            string `mapstructure:"bot_token"`
	ButtonText          string `mapstructure:"button_text"`
	WelcomeMessage      string `mapstructure:"welcome_message"`
	AfterSuccessMessage string `mapstructure:"after_success_message"`
	AfterFailMessage    string `mapstructure:"after_fail_message"`
	PrintSuccessAndFail string `mapstructure:"print_success_and_fail_messages_strategy"`
	WelcomeTimeout      string `mapstructure:"welcome_timeout"`
	BanDurations        string `mapstructure:"ban_duration"`
	UseSocks5Proxy      string `mapstructure:"use_socks5_proxy"`
	Socks5Address       string `mapstructure:"socks5_address"`
	Socks5Port          string `mapstructure:"socks5_port"`
	Socks5Login         string `mapstructure:"socks5_login"`
	Socks5Password      string `mapstructure:"socks5_password"`
}

var config Config

func NewConfig() (Config, error) {
	err := readConfig()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

const CONFIGPATH = "./configs/config.toml"

func readConfig() (err error) {
	v := viper.New()
	path, ok := os.LookupEnv(CONFIGPATH)
	if ok {
		v.SetConfigName("config")
		v.AddConfigPath(path)
	}
	v.SetConfigName("config")
	v.AddConfigPath(".")

	if err = v.ReadInConfig(); err != nil {
		return err
	}
	if err = v.Unmarshal(&config); err != nil {
		return err
	}
	return
}

func GetToken(key string) (string, error) {
	token, ok := os.LookupEnv(key)
	if !ok {
		err := errors.Errorf("Env variable %v isn't set!", key)
		return "", err
	}
	match, err := regexp.MatchString(`^[0-9]+:.*$`, token)
	if err != nil {
		return "", err
	}
	if !match {
		err := errors.Errorf("Telegram Bot Token [%v] is incorrect. Token doesn't comply with regexp: `^[0-9]+:.*$`. Please, provide a correct Telegram Bot Token through env variable TGTOKEN", token)
		return "", err
	}
	return token, nil
}
