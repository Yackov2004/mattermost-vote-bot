package handlers

import (
	"context"
	"fmt"
	"mattermost-voting-bot/internal/settings"
	"mattermost-voting-bot/internal/storage"
	"regexp"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

// handleCreate - хэндлер создания голосования
func handleCreate(application *settings.Application, post *model.Post) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.Logger.Info().
		Str("userID", post.UserId).
		Str("message", post.Message).
		Msg("Получена команда /poll create")

	re := regexp.MustCompile(`"(.*?)"`)
	matches := re.FindAllStringSubmatch(post.Message, -1)
	if len(matches) < 2 {
		application.Logger.Warn().
			Str("userID", post.UserId).
			Msg("Неверный формат команды /poll create")

		SendMsg(application, "Корректный формат запроса: /poll create \"Question\" \"Option1\" \"Option2\"", post.Id)
		return
	}

	question := matches[0][1]
	options := make([]string, 0)
	for _, m := range matches[1:] {
		options = append(options, m[1])
	}
	if len(options) < 2 {
		application.Logger.Warn().
			Str("userID", post.UserId).
			Int("optionsCount", len(options)).
			Msg("Слишком мало вариантов ответа")

		SendMsg(application, "У голосования должно быть как минимум два варианта ответа", post.Id)
		return
	}

	newPoll := &storage.Poll{
		Question: question,
		Options:  options,
		Active:   true,
		OwnerID:  post.UserId,
	}

	id, err := application.Storage.CreatePoll(ctx, newPoll)
	if err != nil {
		application.Logger.Error().
			Err(err).
			Str("userID", post.UserId).
			Msg("Не удалось сохранить голосование в БД")

		SendMsg(application, "Не удалось сохранить голосование в БД", post.Id)
		return
	}

	application.Logger.Info().
		Str("userID", post.UserId).
		Uint64("pollID", id).
		Msg("Голосование успешно создано")

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Создано голосование ID=%d\n", id))
	builder.WriteString(fmt.Sprintf("Вопрос: %s\n", question))
	builder.WriteString("Варианты:\n")
	for i, opt := range options {
		builder.WriteString(fmt.Sprintf("%d) %s\n", i+1, opt))
	}

	SendMsg(application, builder.String(), post.Id)
}
