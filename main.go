package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	_ "github.com/joho/godotenv/autoload"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "drone-dispatch"
	app.Usage = "send a dispatch event to a repository"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "api key to access github api",
			EnvVar: "PLUGIN_API_KEY,GITHUB_RELEASE_API_KEY,GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "dispatch-repo",
			Usage:  "The name of the repo to send the dispatch event to",
			EnvVar: "PLUGIN_DISPATCH_REPO",
		},
		cli.StringFlag{
			Name:   "dispatch-owner",
			Usage:  "The name of the repo owner to send the dispatch event to",
			EnvVar: "PLUGIN_DISPATCH_OWNER",
		},
		cli.StringFlag{
			Name:   "client-data",
			Usage:  "Client data to send with the dispatch event",
			EnvVar: "PLUGIN_CLIENT_DATA",
		},
		cli.StringFlag{
			Name:   "event-type",
			Usage:  "The dispatch request event type",
			EnvVar: "PLUGIN_EVENT_TYPE",
		},
		cli.StringFlag{
			Name:   "repo.fullname",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Event: c.String("build.event"),
		},
		Commit: Commit{
			Ref: c.String("commit.ref"),
		},
		Config: Config{
			APIKey:          c.String("api-key"),
			ClientData:      c.String("client-data"),
			DispatchRepo:    c.String("dispatch-repo"),
			DispatchOwner:	 c.String("dispatch-owner"),
			EventType:		 c.String("event-type"),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
