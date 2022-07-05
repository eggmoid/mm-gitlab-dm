package config

import (
	"errors"
	"strings"
	"sync/atomic"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

var (
	config     atomic.Value
	Mattermost plugin.API
	BotUserID  string
)

type Configuration struct {
	Secret string `json:"secret"`
}

func GetConfig() *Configuration {
	return config.Load().(*Configuration)
}

func SetConfig(c *Configuration) {
	config.Store(c)
}

func (c *Configuration) ProcessConfiguration() error {
	c.Secret = strings.TrimSpace(c.Secret)

	return nil
}

func (c *Configuration) IsValid() error {
	if c.Secret == "" {
		return errors.New("please provide the Webhook Secret")
	}

	return nil
}
