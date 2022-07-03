package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/plugin"
)

// HelloWorldPlugin implements the interface expected by the Mattermost server to communicate
// between the server and plugin processes.
type GitDMPlugin struct {
	plugin.MattermostPlugin
}

type MMDM struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *GitDMPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	m := MMDM{"메시지 페이로드", "mmuser", "general channel"}
	pbytes, _ := json.Marshal(m)
	buff := bytes.NewBuffer(pbytes)
	resp, err := http.Post("http://goq.kro.kr:8080/hooks", "application/json", buff)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Fprint(w, "Hello, world!")
}

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	plugin.ClientMain(&GitDMPlugin{})
}
