package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/plugin"
)

const (
	MMAPI   = ""
	MMTOKEN = ""
	MMBOTID = ""
)

type dict map[string]interface{}

func (d dict) d(k string) dict {
	return d[k].(map[string]interface{})
}

func (d dict) s(k string) string {
	return d[k].(string)
}

type Plugin struct {
	plugin.MattermostPlugin
}

// ServeHTTP demonstrates a plugin that handles HTTP requests
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
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
	url := data.d("object_attributes").s("url")
	title := data.d("object_attributes").s("title")
	description := data.d("object_attributes").s("description")

	// BODY
	// Get username
	username := "admin" // ASSIGNEE
	payload := `author: ` + author + `, url: ` + url + `, title: ` + title + `, description: ` + description
	// BODY

	// Get user id
	userID := getUserID(username)

	// Get channel id
	channelID := getChannelID(MMBOTID, userID)

	// Post message to channel
	createPost(channelID, payload)
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
func createPost(channelID, message string) string {
	var jsonData = []byte(`{
		"channel_id": "` + channelID + `",
		"message": "` + message + `"
	}`)
	req, err := http.NewRequest("POST", MMAPI+"/posts", bytes.NewBuffer(jsonData))
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

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	plugin.ClientMain(&Plugin{})
}
