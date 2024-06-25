package main

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

var ErrorNoTopic = fmt.Errorf("message without topic!")

func get_or_extract_chat(message *tele.Message) int64 {
	var err error
	topic_id := message.ThreadID
	if topic_id == 0 {
		return 0
	}
	chat_id := ChatMap.GetUserChat(topic_id)
	if chat_id == 0 {
		if len(message.ReplyTo.Entities) > 0 {
			user_id := message.ReplyTo.EntityText(message.ReplyTo.Entities[0])
			if user_id != "" {
				chat_id, err = strconv.ParseInt(user_id, 10, 64)
				if err == nil {
					log.Info().Msgf("Reconstructed ChatMap topic_id=%d, chat_id=%d", topic_id, chat_id)
					ChatMap.Pair(chat_id, topic_id)
				} else {
					log.Error().Err(err).Msgf("converting to int user_id=%s", user_id)
				}
			}
		}
	}
	return chat_id
}

func chat_admin(c tele.Context, addReact *tele.ReactionOptions) error {
	var err error
	message := c.Message()
	chat_id := get_or_extract_chat(message)
	if chat_id == 0 {
		return c.Reply(СonfMsg.NoReplyId, tele.ModeHTML)
	}

	var new_message *tele.Message
	var user_chat = &DummyChat{ID: strconv.FormatInt(chat_id, 10)}
	new_message, err = Bot.Copy(user_chat, message)
	if err != nil {
		log.Error().Err(err).Msg("admin_chat: copyMessage to user_chat")
		return err
	}

	// Если указали дополнительную реакцию, применяем ее к отправленному сообщению
	if addReact != nil {
		err = Bot.React(user_chat, &tele.StoredMessage{MessageID: strconv.Itoa(new_message.ID)}, *addReact)
		if err != nil {
			log.Error().Err(err).Msg("admin_chat: addReact to user_chat")
		}
	}

	return nil
}
