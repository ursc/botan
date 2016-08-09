package main

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Application struct {
	Config    *Config
	Questions Questions
	Dict      map[string]string

	Api *telegram.API
	Bot *telebot.Bot
}

func NewApplication(cfgFile, dataFile, dictFile string) *Application {
	cfg := readConfig(cfgFile)
	questions := readQuestions(dataFile)
	dict := readDictionary(dictFile)

	api := telegram.New(cfg.Bot.Token)
	api.Debug(cfg.Bot.Debug)

	bot := telebot.NewWithAPI(api)

	app := &Application{
		Config:    cfg,
		Questions: questions,
		Dict:      dict,
		Api:       api,
		Bot:       bot,
	}

	bot.HandleFunc(app.HandleFunc)
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(app.StartCommand),
		"admin": telebot.CommandFunc(app.AdminCommand),
	}))

	return app
}

func (app *Application) StartCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	msg := newMessage(update.Chat().ID, app.Questions[0])
	_, err := app.Api.SendMessage(ctx, msg)
	return err
}

func (app *Application) AdminCommand(ctx context.Context, arg string) error {
	update := telebot.GetUpdate(ctx)
	chatID := update.Chat().ID
	if chatID != app.Config.Bot.AdminChatId {
		return app.StartCommand(ctx, arg)
	}

	// TODO: ...

	return nil
}

func (app *Application) CallbackFunc(ctx context.Context) error {
	update := telebot.GetUpdate(ctx)

	data := update.CallbackQuery.Data
	if i := strings.IndexByte(data, ':'); i > 0 {
		if data[i+1:] == data[:i] {
			return nil
		}
		data = data[i+1:]
	}

	i, err := strconv.Atoi(data)
	if err != nil {
		return err
	}

	msg := editMessage(update, app.Questions[i])
	_, err = app.Api.EditMessageText(ctx, msg)

	return err
}

func (app *Application) HandleFunc(ctx context.Context) error {
	update := telebot.GetUpdate(ctx)

	if update.CallbackQuery != nil {
		return app.CallbackFunc(ctx)
	}

	var err error
	chatID := update.Chat().ID
	word := strings.ToLower(strings.TrimSpace(update.Message.Text))

	if v, ok := app.Dict[word]; ok {
		msg := telegram.NewMessage(chatID, v)
		_, err = app.Api.SendMessage(ctx, msg)
	} else if chatID > 0 {
		msg := newMessage(chatID, app.Questions[0])
		_, err = app.Api.SendMessage(ctx, msg)
	}

	return err
}

func (app *Application) Start() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !app.Config.Bot.UseWebHook {
		if err = app.Api.SetWebhook(ctx, telegram.NewWebhook("")); err != nil {
			log.Panic(err)
		}
		if err = app.Bot.Serve(ctx); err != nil {
			log.Panic(err)
		}
		return
	}

	webHook := telegram.NewWebhook(app.Config.Bot.WebHook.Host)
	if err = app.Api.SetWebhook(ctx, webHook); err != nil {
		log.Panic(err)
	}
	var h http.HandlerFunc
	if h, err = app.Bot.ServeByWebhook(ctx); err != nil {
		log.Panic(err)
	}

	wc := &app.Config.Bot.WebHook
	err = http.ListenAndServeTLS(wc.Port, wc.CertFile, wc.KeyFile, h)
	if err != nil {
		log.Panic(err)
	}
}

func newMessage(chatID int64, item *Question) telegram.MessageCfg {
	msg := telegram.NewMessage(chatID, item.Title)
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = item.ReplyMarkup

	return msg
}

func editMessage(update *telegram.Update, item *Question) telegram.EditMessageTextCfg {
	msg := telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		item.Title,
	)
	if len(item.Answer) > 0 {
		msg.Text = item.Answer
	}
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = item.ReplyMarkup

	return msg
}
