package main

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func _ban_core(c tele.Context, ban_name string, ban_list map[int64]struct{}) error {
	var err error
	// Команду принимаем только в админском чате
	if c.Chat().ID != Сonf.AdminChat {
		return nil
	}
	message := c.Message()
	// Извлекаем chat_id
	chat_id := get_or_extract_chat(message)
	if chat_id == 0 {
		return c.Reply(СonfMsg.NoReplyId, tele.ModeHTML)
	}

	var ok bool
	_, ok = ban_list[chat_id]
	if !ok {
		ban_list[chat_id] = struct{}{}
		// Оповещаем админа
		err = c.Reply(fmt.Sprintf(СonfMsg.AdmBanned, ban_name), tele.ModeHTML)
		if err != nil {
			log.Error().Err(err).Msg("Cannot send admin MsgBanned")
		}
		if ban_name == "блокировки" {
			// Оповещаем пользователя
			_, err = Bot.Send(&DummyChat{ID: strconv.FormatInt(chat_id, 10)}, СonfMsg.UsrBanned, tele.ModeHTML)
			if err != nil {
				log.Error().Err(err).Msg("Cannot send user MsgBanned")
			}
		}
	} else {
		err = c.Reply(fmt.Sprintf(СonfMsg.AdmBannedAlready, ban_name), tele.ModeHTML)
		if err != nil {
			log.Error().Err(err).Msg("Cannot send admin MsgBannedAlready")
		}
	}
	return nil
}

func InitCommandsMenu() {
	var err error
	// User
	err = Bot.SetCommands(
		tele.CommandScope{Type: tele.CommandScopeDefault},
		[]tele.Command{
			{Text: "help", Description: "Информация о боте"},
		},
	)
	// Admin
	err = Bot.SetCommands(
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: Сonf.AdminChat},
		[]tele.Command{
			{Text: "ban", Description: "Заблокировать пользователя (с сообщением)"},
			{Text: "ignore", Description: "Игнорировать пользователя (без сообщений)"},
			{Text: "unblock", Description: "Снять все блокировки с пользователя"},
			// {Text: "restrict_list", Description: "Показать список заблокированных и игнорируемых"},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("bot SetCommands")
	}

	Bot.Handle("/help", func(c tele.Context) error {
		return c.Send(СonfMsg.Help, tele.ModeHTML, tele.NoPreview)
	})

	Bot.Handle("/start", func(c tele.Context) error {
		if Сonf.StartMsg != "" {
			c.Send(Сonf.StartMsg, tele.ModeHTML)
		}
		return c.Send(СonfMsg.Help, tele.ModeHTML, tele.NoPreview)
	})

	Bot.Handle("/ban", func(c tele.Context) error {
		return _ban_core(c, "блокировки", ListBan)
	})

	Bot.Handle("/ignore", func(c tele.Context) error {
		return _ban_core(c, "игнорирования", ListIgnore)
	})

	Bot.Handle("/unblock", func(c tele.Context) error {
		// Команду принимаем только в админском чате
		if c.Chat().ID != Сonf.AdminChat {
			return nil
		}
		message := c.Message()

		// Извлекаем chat_id
		chat_id := get_or_extract_chat(message)
		if chat_id == 0 {
			return c.Reply(СonfMsg.NoReplyId, tele.ModeHTML)
		}
		delete(ListBan, chat_id)
		delete(ListIgnore, chat_id)

		c.Reply(СonfMsg.AdmUnBanned, tele.ModeHTML)
		return nil
	})
}
