package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eggmoid/mm-gitlab-dm/config"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	MMDOMAIN = ""
	MMAPI    = ""
	MMTOKEN  = ""
	MMBOTID  = ""
)

type dict map[string]interface{}

func (d dict) d(k string) dict {
	return d[k].(map[string]interface{})
}

func (d dict) s(k string) string {
	return d[k].(string)
}

type GitPlugin struct {
	plugin.MattermostPlugin
}

func (p *GitPlugin) OnActivate() error {
	config.Mattermost = p.API

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	return nil
}

func (p *GitPlugin) OnConfigurationChange() error {
	if config.Mattermost == nil {
		return nil
	}
	var configuration config.Configuration

	config.SetConfig(&configuration)
	return nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests
func (p *GitPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var data dict
	json.Unmarshal([]byte(body), &data)

	if data.s("event_type") != "merge_request" {
		return
	}

	author := data.d("user").s("username")
	name := data.d("user").s("name")
	url := data.d("object_attributes").s("url")
	title := data.d("object_attributes").s("title")
	description := data.d("object_attributes").s("description")
	namespace := data.d("project").s("namespace")
	project := data.d("project").s("name")
	project_url := data.d("project").s("homepage")

	for _, a := range data["assignees"].([]interface{}) {
		username := a.(map[string]interface{})["username"].(string)
		payload := name + ` (` + author + `) opened merge request ` + `[` + title + `](` + url + `) in [` + namespace + ` / ` + project + `](` + project_url + `)`

		// Get user id
		userID := getUserID(username)

		// Get channel id
		channelID := getChannelID(MMBOTID, userID)

		// Post message to channel
		createPost(channelID, payload, title, url, description)
	}
}

// get string and return string {{{
func getUserID(username string) string {
	req, err := http.NewRequest("GET", MMAPI+"/users/username/"+username, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+MMTOKEN)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)
	data := map[string]string{}

	json.Unmarshal([]byte(str), &data)
	return data["id"]
}

// }}}

// get channel id from two user id {{{
func getChannelID(botID, userID string) string {
	var jsonData = []byte(`["` + botID + `", "` + userID + `"]`)
	req, err := http.NewRequest("POST", MMAPI+"/channels/direct", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+MMTOKEN)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)
	data := map[string]string{}

	json.Unmarshal([]byte(str), &data)
	return data["id"]
}

// }}}

// create a post {{{
func createPost(channelID, message, title, title_link, text string) string {
	attachment := &model.SlackAttachment{
		Fallback:  "",
		Color:     "#db3b21",
		Title:     title,
		TitleLink: title_link,
		Text:      text,
	}
	post := &model.Post{
		UserId:    MMBOTID,
		ChannelId: channelID,
		Message:   message,
	}

	model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})

	client := model.NewAPIv4Client(MMDOMAIN)
	client.SetToken(MMTOKEN)

	client.CreatePost(post)

	return ""
}

// }}}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	plugin.ClientMain(&GitPlugin{})
}
