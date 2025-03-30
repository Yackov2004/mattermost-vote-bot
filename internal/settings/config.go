package settings

import (
	"net/url"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// Config содержит конфиг бота
type Config struct {
	MattermostUserName string
	MattermostTeamName string
	MattermostToken    string
	MattermostChannel  string
	MattermostServer   *url.URL
}

// LoadConfig загружает в конфиг переменные среды
func LoadConfig() Config {
	var settings Config

	settings.MattermostTeamName = os.Getenv("MM_TEAM")
	settings.MattermostUserName = os.Getenv("MM_USERNAME")
	settings.MattermostToken = os.Getenv("MM_TOKEN")
	settings.MattermostChannel = os.Getenv("MM_CHANNEL")
	settings.MattermostServer, _ = url.Parse(os.Getenv("MM_SERVER"))

	return settings
}
