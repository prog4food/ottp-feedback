package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func get_or_create_topic(message *tele.Message, force bool) int {
	sender_id := message.Sender.ID
	topic_id := ChatMap.GetMsgTopic(sender_id)
	if topic_id == 0 || force == true {
		topic_msg, err := Bot.Send(
			Chat,
			fmt.Sprintf("\n💬 <code>%d</code>\nuser: @%s\nname: <code>%s %s</code>", sender_id, message.Sender.Username, message.Sender.FirstName, message.Sender.LastName),
			tele.ModeHTML,
		)
		if err != nil {
			log.Error().Err(err).Msg("sending topic message")
		} else {
			topic_id = topic_msg.ID
			ChatMap.Pair(sender_id, topic_id)
		}
	}
	return topic_id
}

func chat_user(c tele.Context, addReact *tele.ReactionOptions) error {
	var err error
	message := c.Message()
	user_id := message.Sender.ID

	var ok bool
	// Игнор лист
	_, ok = ListIgnore[user_id]
	if ok {
		return nil
	}
	// Бан лист
	_, ok = ListBan[user_id]
	if ok {
		return c.Send(СonfMsg.UsrBanned, tele.ModeHTML)
	}

	var new_message *tele.Message
	// Открываем топик, или создаем его
	topic_id := get_or_create_topic(message, false)
	// Пересылаем сообщение в админиский чат
	new_message, err = Bot.Copy(Chat, message, &tele.SendOptions{
		ReplyTo: &tele.Message{ID: topic_id},
	})
	if err != nil && strings.Contains(err.Error(), "message to be replied not found") {
		log.Warn().Msg("message to be replied not found, try to open new topic")
		// Принудительно создаем топик, и пробуем еще раз переслать сообщение
		topic_id = get_or_create_topic(message, true)
		new_message, err = Bot.Copy(Chat, message, &tele.SendOptions{
			ReplyTo: &tele.Message{ID: topic_id},
			// AllowWithoutReply: true,
		})
	}
	if err != nil {
		log.Error().Err(err).Msg("user_chat: copyMessage to admin_chat")
		return err
	}

	// Если указали дополнительную реакцию, применяем ее к отправленному сообщению
	if addReact != nil {
		err = Bot.React(Chat, &tele.StoredMessage{MessageID: strconv.Itoa(new_message.ID)}, *addReact)
		if err != nil {
			log.Error().Err(err).Msg("user_chat: addReact to admin_chat")
		}
	}

	botNotes := make([]string, 0, 3)
	if message.LastEdit == 0 {
		// Если сообщение новое (LastEdit == 0)
		//  добавляем на него реакцию reactSended
		err = Bot.React(message.Chat, message, reactSended)
		if err != nil {
			log.Error().Err(err).Msg("user_chat: reactSended")
		}
	} else {
		// Если сообщение отредактировано (LastEdit != 0)
		//  сообщаем, что это немного не так работает
		botNotes = append(botNotes, СonfMsg.WarnEdit)
	}

	if message.ReplyTo != nil {
		// Если пользователь ответил на сообщение
		//  сообщаем, что это немного не так работает
		botNotes = append(botNotes, СonfMsg.WarnReply)
	}

	// botNote - заметки от бота
	if len(botNotes) > 0 {
		err = c.Reply(strings.Join(botNotes, "\n\n"))
		if err != nil {
			log.Error().Err(err).Msg("user_chat: botNotes")
		}
	}
	return nil
}
