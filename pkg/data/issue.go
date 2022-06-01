package data

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v45/github"
	stripmd "github.com/writeas/go-strip-markdown"
	"golang.org/x/oauth2"
)

type Repo struct {
	Name        string
	Issues      []Entry
	Discussions []Entry
}

func FromGHIssues(repos []string) (Set, error) {
	var input Set

	issues, err := getQuestions(repos)
	if err != nil {
		return input, err
	}
	input.Mapping = Mapping{
		Path: "",
		Fields: map[string]string{
			"title":  "string",
			"text":   "string",
			"author": "string",
			"date":   "numeric",
			"type":   "string",
		},
	}

	sources, err := readIssues(issues)
	if err != nil {
		return input, err
	}
	input.Sources = sources

	return input, nil
}

func readIssues(entries []Repo) ([]Source, error) {
	var sources []Source

	for _, entry := range entries {
		for _, issue := range entry.Issues {
			_id := fmt.Sprintf(
				"title=%s&url=%s&author=%s&year=%d&type=%s",
				issue.Title, issue.URL, issue.Author, issue.Year, "issue")

			sources = append(sources, Source{
				ID:  _id,
				URL: issue.URL,
				Fields: map[string]interface{}{
					"title":  issue.Title,
					"text":   issue.Body,
					"author": issue.Author,
					"year":   issue.Year,
					"type":   issue.Type,
				},
			})
		}
	}

	return sources, nil
}

func getQuestions(repos []string) ([]Repo, error) {
	var loaded []Repo

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	for _, name := range repos {
		repo := Repo{Name: "Vale"}

		_, resp, err := client.Issues.ListByRepo(ctx, "errata-ai", name, &github.IssueListByRepoOptions{State: "all"})
		if err != nil {
			return loaded, err
		}

		fmt.Println(fmt.Sprintf("There are %d pages.", resp.LastPage))
		for i := 1; i <= resp.LastPage; i++ {
			issues, _, err := client.Issues.ListByRepo(
				ctx, "errata-ai", name, &github.IssueListByRepoOptions{
					State:       "all",
					ListOptions: github.ListOptions{Page: i},
				})

			if err != nil {
				return loaded, err
			}

			for _, issue := range issues {
				year, _, _ := issue.CreatedAt.Date()
				ent := Entry{
					Title: *issue.Title,
					URL:   *issue.HTMLURL,
					Type:  "issue",
					Year:  year,
				}
				if issue.GetUser().Login != nil {
					ent.Author = *issue.GetUser().Login
				}
				if issue.Body != nil {
					ent.Body = stripmd.Strip(*issue.Body)
				}
				repo.Issues = append(repo.Issues, ent)
			}
		}

		loaded = append(loaded, repo)
	}

	return loaded, nil
}
