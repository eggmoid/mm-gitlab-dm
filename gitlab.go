package main

import (
	"log"

	fn "github.com/thoas/go-funk"
	"github.com/xanzy/go-gitlab"
)

func retrieveUsernames(id, merge_request_iid int) ([]string, error) {
	git, err := gitlab.NewClient(GITLAB_TOKEN, gitlab.WithBaseURL(GITLAB_URL))
	if err != nil {
		return nil, err
	}

	notes, _, err := git.Notes.ListMergeRequestNotes(id, merge_request_iid, &gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		log.Printf("Failed to fetch notes: %v", err)
		return nil, err
	}

	return fn.Uniq(fn.Map(notes, func(note *gitlab.Note) string {
		return note.Author.Username
	})).([]string), nil
}
