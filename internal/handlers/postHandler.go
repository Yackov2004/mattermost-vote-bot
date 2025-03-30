package handlers

import (
	"mattermost-voting-bot/internal/settings"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

// HandlePost - хэндлер для роутинга запросов
func HandlePost(application *settings.Application, post *model.Post) {
	if !strings.HasPrefix(post.Message, "/poll") {
		return
	}

	args := strings.Fields(post.Message)
	if len(args) < 2 {
		SendMsg(application, "Корректный формат запроса: /poll create/vote/results/close/delete", post.Id)
		return
	}

	action := strings.ToLower(args[1])
	switch action {
	case "create":
		handleCreate(application, post)
	case "vote":
		handleVote(application, post)
	case "results":
		handleResults(application, post)
	case "close":
		handleClose(application, post)
	case "delete":
		handleDelete(application, post)
	default:
		SendMsg(application, "Недопустимая команда. Список команд: create, vote, results, close, delete", post.Id)
	}
}
