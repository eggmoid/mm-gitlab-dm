package main

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
	fn "github.com/thoas/go-funk"
	"golang.org/x/oauth2"
)

const (
	GITLAB_TOKEN = ""
	GITLAB_URL   = ""
)

func retrieveUserIDsByMRID(id string) ([]string, error) {
	httpClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GITLAB_TOKEN},
		))
	client := graphql.NewClient(GITLAB_URL, httpClient)
	type MergeRequestID string

	type MRNode struct {
		Author struct {
			ID string `graphql:"id"`
		} `graphql:"author"`
	}

	var query struct {
		MergeRequest struct {
			Notes struct {
				Nodes []struct {
					Author struct {
						ID string `graphql:"id"`
					} `graphql:"author"`
				} `graphql:"nodes"`
			} `graphql:"notes"`
		} `graphql:"mergeRequest(id: $mrId)"`
	}
	variables := map[string]interface{}{
		"mrId": MergeRequestID("gid://gitlab/MergeRequest/" + id),
	}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	ids := fn.Uniq(fn.Map(query.MergeRequest.Notes.Nodes, func(n MRNode) string {
		return n.Author.ID
	}))
	return ids.([]string), nil
}

func retrieveUsernamesByUserID(ids []string) ([]string, error) {
	httpClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GITLAB_TOKEN},
		))
	client := graphql.NewClient(GITLAB_URL, httpClient)

	type UserIDs []string
	type UserNode struct {
		Nodes []struct {
			Username string `graphql:"username"`
		} `graphql:"nodes"`
	}

	var query struct {
		Users struct {
			UserNode
		} `graphql:"users(ids: $ids)"`
	}

	variables := map[string]interface{}{
		"ids": UserIDs(ids),
	}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	return fn.Map(query.Users.UserNode.Nodes, func(n struct {
		Username string `graphql:"username"`
	}) string {
		return n.Username
	}).([]string), nil
}
