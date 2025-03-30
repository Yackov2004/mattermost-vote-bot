package app

import (
	"encoding/json"
	"fmt"
	"mattermost-voting-bot/internal/handlers"
	"mattermost-voting-bot/internal/settings"
	"mattermost-voting-bot/internal/storage"
	"os"
	"os/signal"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/rs/zerolog"
)

// Run - запуск приложения
func Run() {
	app := &settings.Application{
		Logger: zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC822,
			},
		).With().Timestamp().Logger(),
	}

	app.Config = settings.LoadConfig()
	app.Logger.Info().Str("Config", fmt.Sprint(app.Config)).Msg("")

	setupGracefulShutdown(app)

	// Подключение к Tarantool
	st, err := storage.NewStorage(os.Getenv("TARANTOOL_HOST"), os.Getenv("TARANTOOL_PORT"))
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Не удалось подключиться к БД")
	}
	app.Storage = st

	// Авторизация в Mattermost
	app.MattermostClient = model.NewAPIv4Client(app.Config.MattermostServer.String())
	app.MattermostClient.SetToken(app.Config.MattermostToken)

	user, resp, err := app.MattermostClient.GetUser("me", "")
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Не удалось авторизироваться")
	}
	app.Logger.Debug().Interface("user", user).Interface("resp", resp).Msg("")
	app.Logger.Info().Msg("Авторизация в mattermost успешна")
	app.MattermostUser = user

	team, resp, err := app.MattermostClient.GetTeamByName(app.Config.MattermostTeamName, "")
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Could not find team")
	}
	app.Logger.Debug().Interface("team", team).Interface("resp", resp).Msg("")
	app.MattermostTeam = team

	channel, resp, err := app.MattermostClient.GetChannelByName(
		app.Config.MattermostChannel, app.MattermostTeam.Id, "",
	)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Could not find channel")
	}
	app.Logger.Debug().Interface("channel", channel).Interface("resp", resp).Msg("")
	app.MattermostChannel = channel

	//handlers.SendMsg(app, "Hi! I am a bot.", "")

	listenToEvents(app)
}

func setupGracefulShutdown(app *settings.Application) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if app.MattermostWebSocketClient != nil {
				app.Logger.Info().Msg("Closing websocket connection")
				app.MattermostWebSocketClient.Close()
			}
			app.Logger.Info().Msg("Shutting down")
			os.Exit(0)
		}
	}()
}

func listenToEvents(app *settings.Application) {
	var err error
	failCount := 0
	for {
		app.MattermostWebSocketClient, err = model.NewWebSocketClient4(
			fmt.Sprintf("ws://%s", app.Config.MattermostServer.Host+app.Config.MattermostServer.Path),
			app.MattermostClient.AuthToken,
		)
		if err != nil {
			app.Logger.Warn().Err(err).Msg("Mattermost websocket disconnected, retrying")
			failCount += 1
			time.Sleep(time.Duration(failCount) * time.Second)
			continue
		}
		app.Logger.Info().Msg("Mattermost websocket connected")

		app.MattermostWebSocketClient.Listen()

		for event := range app.MattermostWebSocketClient.EventChannel {
			go handleWebSocketEvent(app, event)
		}
	}
}

func handleWebSocketEvent(app *settings.Application, event *model.WebSocketEvent) {

	if event.GetBroadcast().ChannelId != app.MattermostChannel.Id {
		return
	}

	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	post := &model.Post{}
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		app.Logger.Error().Err(err).Msg("Could not cast event to *model.Post")
		return
	}

	if post.UserId == app.MattermostUser.Id {
		return
	}

	handlers.HandlePost(app, post)
}
