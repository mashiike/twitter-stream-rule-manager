package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

var (
	rulesFilePath string
	bearer        string
)

func main() {
	subcommands.Register(subcommands.CommandsCommand(), "help")
	subcommands.Register(subcommands.FlagsCommand(), "help")
	subcommands.Register(subcommands.HelpCommand(), "help")

	subcommands.Register(&diffCmd{}, "")
	subcommands.Register(&deployCmd{}, "")

	setFlags(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	subcommands.Execute(ctx)
}

func setFlags(f *flag.FlagSet) {
	defaultToken := os.Getenv("TWITTER_BEARER_TOKEN")
	f.StringVar(&rulesFilePath, "rules", "rules.json", "rule file path (required)")
	f.StringVar(&bearer, "bearer", defaultToken, "twitter bearer token (required, or set with TWITTER_BEARER_TOKEN env)")
}
