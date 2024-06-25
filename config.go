package main

import (
	"os"
	"strconv"

	"github.com/hjson/hjson-go/v4"
	"github.com/rs/zerolog/log"
)

type ConfigMsg_t struct {
	Help             string `json:"help"`
	NoReplyId        string `json:"no_reply_id"`
	AdmBanned        string `json:"adm_banned"`
	AdmBannedAlready string `json:"adm_banned_already"`
	AdmUnBanned      string `json:"adm_unbanned"`
	UsrBanned        string `json:"usr_banned"`
	WarnEdit         string `json:"warn_edit"`
	WarnReply        string `json:"warn_reply"`
}

type Config_t struct {
	AdminChat     int64  `json:"admin_chat"`
	BotToken      string `json:"bot_token"`
	WebhookListen string `json:"webhook_listen_local"`
	WebhookUrl    string `json:"webhook_public_url"`
	StartMsg      string `json:"start_message"`
}

const (
	fileConfig    = "config.hjson"
	fileConfigMsg = "messages.hjson"
)

var (
	Сonf    *Config_t
	СonfMsg *ConfigMsg_t
)

func read_config() {
	Сonf = &Config_t{}
	jsonData, err := os.ReadFile(fileConfig)
	if err == nil {
		err = hjson.Unmarshal(jsonData, Сonf)
	}
	if err != nil {
		log.Err(err).Msg(fileConfig + ": read config")
	}

	// Грузим BotToken из переменной окружения, если она есть
	enVar := os.Getenv("BOT_TOKEN")
	if enVar != "" {
		Сonf.BotToken = enVar
	}

	if Сonf.BotToken == "" {
		log.Fatal().Msg("BotToken - is not set!")
	}

	// Грузим AdminChat из переменной окружения, если она есть
	enVar = os.Getenv("ADMIN_CHAT")
	if enVar != "" {
		Сonf.AdminChat, err = strconv.ParseInt(enVar, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("parse enviroment ADMIN_CHAT")
		}
	}

	// Грузим StartMsg из переменной окружения, если она есть
	enVar = os.Getenv("START_MESSAGE")
	if enVar != "" {
		Сonf.StartMsg = enVar
	}

	СonfMsg = &ConfigMsg_t{}
	jsonData, err = os.ReadFile(fileConfigMsg)
	if err == nil {
		err = hjson.Unmarshal(jsonData, СonfMsg)
	}
	if err != nil {
		log.Err(err).Msg(fileConfig + ": cannot read config!!!")
	}
	// (с) //
	СonfMsg.Help += "\n\n" + `<a href="https://github.com/prog4food">prog4food</a>`
}
