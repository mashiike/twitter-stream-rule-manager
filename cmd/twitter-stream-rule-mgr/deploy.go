package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
	twstrulemgr "github.com/mashiike/twitter-filterd-stream-rule-manager"
)

type deployCmd struct {
	dryRun bool
}

func (c *deployCmd) Name() string { return "deploy" }

func (c *deployCmd) Synopsis() string { return "deploy" }

func (c *deployCmd) Usage() string { return "deploy" }

func (c *deployCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&rulesFilePath, "rules", rulesFilePath, "rule file ")
	f.StringVar(&bearer, "bearer", "", "twitter bearer token (required)")
	f.BoolVar(&c.dryRun, "dry-run", false, "dry run flag")
}

func (c *deployCmd) Execute(ctx context.Context, _ *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	config := &twstrulemgr.Config{
		BearerToken: bearer,
		RulesFile:   rulesFilePath,
	}
	app, err := twstrulemgr.New(config)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	if err := app.Deploy(ctx, c.dryRun); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
