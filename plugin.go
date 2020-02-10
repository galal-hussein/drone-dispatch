package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/v29/github"
)

type (
	Repo struct {
		Owner   string
		Name    string
	}

	Build struct {
		Event    string
	}

	Commit struct {
		Ref     string
	}
	Config struct {
		APIKey          string
		TargetRepo		string
		TargetOwner     string
		ClientData		string
		EventType		string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	if p.Build.Event != "tag" {
		return fmt.Errorf("The GitHub Dispatch plugin is only available for tags")
	}

	if p.Config.APIKey == "" {
		return fmt.Errorf("You must provide an API key")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Config.APIKey})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)


	jsonClientData := json.RawMessage(p.Config.ClientData)
	dispatchOptions := github.DispatchRequestOptions{
		EventType:     p.Config.EventType,
		ClientPayload: &jsonClientData,
	}


	_, _, err := client.Repositories.Dispatch(ctx, p.Config.TargetOwner, p.Config.TargetRepo, dispatchOptions)
	if err != nil {
		return fmt.Errorf("Failed to send dispatch to repo %s: %v", p.Config.TargetRepo, err)
	}
	return nil
}
