package handlers

import (
	"mattermost-voting-bot/internal/settings"

	"github.com/mattermost/mattermost-server/v6/model"
)

// SendMsg - функция для отправки сообщения
func SendMsg(application *settings.Application, msg string, replyToId string) {
	post := &model.Post{
		ChannelId: application.MattermostChannel.Id,
		Message:   msg,
	}

	if _, _, err := application.MattermostClient.CreatePost(post); err != nil {
		application.Logger.Error().Err(err).Str("Msg Id", replyToId).Msg("Не удалось послать сообщение")
	}
}
