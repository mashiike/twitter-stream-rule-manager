package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
	twstrulemgr "github.com/mashiike/twitter-stream-rule-manager"
)

type diffCmd struct{}

func (c *diffCmd) Name() string { return "diff" }

func (c *diffCmd) Synopsis() string { return "diff" }

func (c *diffCmd) Usage() string { return "diff" }

func (c *diffCmd) SetFlags(f *flag.FlagSet) {
	setFlags(f)
}

func (c *diffCmd) Execute(ctx context.Context, _ *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	config := &twstrulemgr.Config{
		BearerToken: bearer,
		RulesFile:   rulesFilePath,
	}
	app, err := twstrulemgr.New(config)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	if err := app.Diff(ctx); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
