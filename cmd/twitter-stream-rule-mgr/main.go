package main

import (
	"context"
	"flag"

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

	flag.StringVar(&rulesFilePath, "rules", "rules.json", "rule file path (required)")
	flag.StringVar(&bearer, "bearer", "", "twitter bearer token (required)")
	flag.Parse()

	ctx := context.Background()
	subcommands.Execute(ctx)
}
