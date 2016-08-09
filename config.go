package main

import (
	"github.com/BurntSushi/toml"
	"net/url"
	"strings"
)

type Config struct {
	Bot struct {
		Token       string
		UseWebHook  bool
		AdminChatId int64
		Debug       bool

		WebHook struct {
			Host     string
			Port     string
			CertFile string
			KeyFile  string
		}
	}
}

func readConfig(path string) *Config {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	checkErr(err)

	u, err := url.Parse(cfg.Bot.WebHook.Host)
	checkErr(err)

	if index := strings.LastIndexByte(u.Host, ':'); index != -1 {
		cfg.Bot.WebHook.Port = u.Host[index:]
	} else {
		cfg.Bot.WebHook.Port = ":https"
	}

	return &cfg
}
