package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/react"
)

var (
	Bot        *tele.Bot
	BotId      int64
	Chat       *DummyChat
	ChatMap    *ChatMap_t
	ListBan    map[int64]struct{} = make(map[int64]struct{})
	ListIgnore map[int64]struct{} = make(map[int64]struct{})

	reactEdited tele.ReactionOptions = react.React(tele.Reaction{Type: "emoji", Emoji: "✍"})
	reactSended tele.ReactionOptions = react.React(tele.Reaction{Type: "emoji", Emoji: "🕊"})
)

// Устанавливаются при сборке
var depl_ver = "[devel]"

// signal handler
func signalHandler(signal os.Signal) {
	log.Warn().Msgf("Caught signal: %+v", signal)
	isTermSignal := false

	switch signal {

	case syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT:
		isTermSignal = true
	default:
		log.Error().Msgf("Unknown signal: %+v", signal)
	}

	if isTermSignal {
		log.Info().Msg("Shutdown...")
		os.Exit(0)
	}
}

// initialize signal handler
func initSignals() {
	captureSignal := make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	signalHandler(<-captureSignal)
}

func main() {
	var err error
	InitLogger("OTT-play FOSS feeedback bot", depl_ver, false)
	read_config()

	// Создаем глобальные объекты
	Chat = &DummyChat{ID: strconv.FormatInt(Сonf.AdminChat, 10)}
	ChatMap = NewChatMap()

	var pooler tele.Poller
	allowedupdates := []string{"message", "edited_message"}
	if Сonf.WebhookListen != "" && Сonf.WebhookUrl != "" {
		// В конфиге установлены параметры WebHook, используем его
		pooler = &tele.Webhook{
			Listen:         Сonf.WebhookListen,
			Endpoint:       &tele.WebhookEndpoint{PublicURL: Сonf.WebhookUrl},
			AllowedUpdates: allowedupdates,
		}
	} else {
		// Используем стандартный LongPoller
		pooler = &tele.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: allowedupdates,
		}
	}

	pref := tele.Settings{
		Token: Сonf.BotToken, Poller: pooler,
		// Verbose: true,
	}

	Bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	// Если не используется Webhook, то предварительно удаляем его
	switch Bot.Poller.(type) {
	case *tele.Webhook:
		log.Info().Msg("Working in WebHook mode")
	default:
		log.Info().Msg("Working in Poller mode")
		Bot.RemoveWebhook(false)
	}

	// Если AdminChat не установлен, то на все сообщения отвечаем  активным chat_id
	if Сonf.AdminChat == 0 {
		log.Warn().Msg("AdminChat is empty. The bot will respond current chat_id to all messages!")
		Bot.Handle(tele.OnText, func(c tele.Context) error {
			return c.Reply("Current chat_id=<code>"+strconv.FormatInt(c.Chat().ID, 10)+"</code>", tele.ModeHTML)
		})
		Bot.Start()
		return
	}
	go initSignals()

	// Инициализация бота
	BotId = Bot.Me.ID
	InitCommandsMenu()

	_register_std_bot_handle := func(tgAction string, addReact *tele.ReactionOptions) {
		Bot.Handle(tgAction, func(c tele.Context) error {
			// fmt.Printf("%+v\n", c)
			log.Debug().Msgf("message from %d [%s]", c.Message().Sender.ID, tgAction)
			chat := c.Chat()
			if chat.ID == Сonf.AdminChat {
				return chat_admin(c, addReact)
			}
			// Реагируем только на личные сообщения боту
			if chat.Type == tele.ChatPrivate {
				return chat_user(c, addReact)
			}
			return nil
		})
	}

	// Реагируем на:
	_register_std_bot_handle(tele.OnText, nil)      // Обычные сообщения
	_register_std_bot_handle(tele.OnPhoto, nil)     // Фотографии
	_register_std_bot_handle(tele.OnDocument, nil)  // Документы
	_register_std_bot_handle(tele.OnAudio, nil)     // Аудио записи
	_register_std_bot_handle(tele.OnVoice, nil)     // Голосовые сообщения
	_register_std_bot_handle(tele.OnVideo, nil)     // Видео
	_register_std_bot_handle(tele.OnVideoNote, nil) // Кружочки
	_register_std_bot_handle(tele.OnContact, nil)   // Контакты
	_register_std_bot_handle(tele.OnLocation, nil)  // Местоположение
	_register_std_bot_handle(tele.OnAnimation, nil) // Гифки
	_register_std_bot_handle(tele.OnSticker, nil)   // Стикеры
	// _register_std_bot_handle(tele.OnVenue, nil)    // Место встречи (не работает)
	// _register_std_bot_handle(tele.OnPoll, nil)     // Голосование (не работает)

	// Обработчик редактирования сообщений
	_register_std_bot_handle(tele.OnEdited, &reactEdited) // Стикеры

	Bot.Start()
}
