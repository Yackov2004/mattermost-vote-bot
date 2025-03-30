package handlers

import (
	"context"
	"fmt"
	"mattermost-voting-bot/internal/settings"
	"mattermost-voting-bot/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

// handleVote - хэндлер регистрации голоса
func handleVote(app *settings.Application, post *model.Post) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Logger.Info().
		Str("userID", post.UserId).
		Str("message", post.Message).
		Msg("Получена команда /poll vote")

	parts := strings.Fields(post.Message)
	if len(parts) < 4 {
		app.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Неверный формат команды /poll vote")

		SendMsg(app, "Корректный формат запроса: /poll vote <pollID> \"Option\"", post.Id)
		return
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		app.Logger.Warn().
			Str("userID", post.UserId).
			Str("rawPollID", parts[2]).
			Msg("Не удалось распарсить pollID")

		SendMsg(app, "Invalid pollID", post.Id)
		return
	}

	re := regexp.MustCompile(`"(.*?)"`)
	match := re.FindStringSubmatch(post.Message)
	if len(match) < 2 {
		app.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Ошибка в формате запроса")

		SendMsg(app, "Ошибка в формате запроса", post.Id)
		return
	}
	chosenOption := match[1]

	poll, err := app.Storage.GetPoll(ctx, pollID)
	if err != nil {
		app.Logger.Error().
			Err(err).
			Uint64("pollID", pollID).
			Msg("Не удалось найти голосование в БД")

		SendMsg(app, fmt.Sprintf("Не удалось найти голосование: %d", pollID), post.Id)
		return
	}
	if !poll.Active {
		app.Logger.Warn().
			Uint64("pollID", pollID).
			Msg("Голосование уже закрыто")

		SendMsg(app, "Это голосование уже закрыто", post.Id)
		return
	}

	// Проверяем, есть ли такой вариант
	found := false
	for _, opt := range poll.Options {
		if opt == chosenOption {
			found = true
			break
		}
	}
	if !found {
		app.Logger.Warn().
			Str("chosenOption", chosenOption).
			Uint64("pollID", pollID).
			Msg("Нет такого варианта ответа")

		SendMsg(app, "Нет такого варианта ответа", post.Id)
		return
	}

	v := &storage.Vote{
		PollID: pollID,
		UserID: post.UserId,
		Option: chosenOption,
	}

	err = app.Storage.CreateVote(ctx, v)
	if err != nil {
		// Если duplicate key - пользователь уже голосовал
		if strings.Contains(err.Error(), "Duplicate key exists") {
			app.Logger.Info().
				Str("userID", post.UserId).
				Uint64("pollID", pollID).
				Str("option", chosenOption).
				Msg("Пользователь уже голосовал")

			SendMsg(app, "Вы уже голосовали", post.Id)
			return
		}
		app.Logger.Error().
			Err(err).
			Str("userID", post.UserId).
			Uint64("pollID", pollID).
			Str("option", chosenOption).
			Msg("Ошибка записи варианта ответа в БД")

		SendMsg(app, "Ошибка записи варианта ответа в БД", post.Id)
		return
	}

	app.Logger.Info().
		Str("userID", post.UserId).
		Uint64("pollID", pollID).
		Str("option", chosenOption).
		Msg("Голос записан успешно")

	SendMsg(app, fmt.Sprintf("Вы проголосовали за вариант \"%s\"", chosenOption), post.Id)
}
