package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/v29/github"
	"strings"
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
		DispatchRepo	string
		DispatchOwner   string
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
	if p.Config.DispatchOwner == "" {
		return fmt.Errorf("You must provide an Target owner for dispatch repo: %#v", p.Config)
	}
	if p.Config.DispatchRepo == "" {
		return fmt.Errorf("You must provide an Target Repo to dispatch")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Config.APIKey})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	clientData := strings.TrimSuffix(p.Config.ClientData, "\n")
	bytes, err := json.Marshal(clientData)
	if err != nil {
		return fmt.Errorf("Failed to marshal json client data: %v", err)
	}
	payload := json.RawMessage(bytes)
	dispatchOptions := github.DispatchRequestOptions{
		EventType:     p.Config.EventType,
		ClientPayload: &payload,
	}

	_, _, err = client.Repositories.Dispatch(ctx, p.Config.DispatchOwner, p.Config.DispatchRepo, dispatchOptions)
	if err != nil {
		return fmt.Errorf("Failed to send dispatch to repo %s: %v with client_data: %v", p.Config.DispatchRepo, err, string(payload))
	}
	return nil
}
