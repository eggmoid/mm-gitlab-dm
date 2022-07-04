package config

import (
	"sync/atomic"

	"github.com/mattermost/mattermost-server/v6/plugin"
)

var (
	config     atomic.Value
	Mattermost plugin.API
	BotUserID  string
)

type Configuration struct {
	Secret string `json:"secret"`
}
