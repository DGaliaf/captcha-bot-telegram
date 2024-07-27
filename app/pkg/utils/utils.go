package utils

import (
	"captcha-bot/app/internal/config"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
)

func GetBanDuration(cfg config.Config) (int64, error) {
	if cfg.BanDurations == "forever" {
		return tb.Forever(), nil
	}

	n, err := strconv.ParseInt(cfg.BanDurations, 10, 64)
	if err != nil {
		return 0, err
	}

	return time.Now().Add(time.Duration(n) * time.Minute).Unix(), nil
}

func InitSocks5Client(cfg config.Config) (*http.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Socks5Address, cfg.Socks5Port)
	dialer, err := proxy.SOCKS5("tcp", addr, &proxy.Auth{User: cfg.Socks5Login, Password: cfg.Socks5Password}, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("cannot init socks5 proxy client dialer: %w", err)
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	httpTransport.DialContext = dialContext

	return httpClient, nil
}
