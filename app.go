package twstrulemgr

import (
	"context"
	"log"
	"time"
)

type App struct {
	Rules  Rules
	Client *Client
}

func New(config *Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	app := &App{
		Rules:  config.Rules,
		Client: newClient(config),
	}
	return app, nil
}

func (app *App) Diff(ctx context.Context) error {
	_, err := app.diff(ctx)
	return err
}

func (app *App) diff(ctx context.Context) (DiffRules, error) {
	rules, sent, err := app.Client.GetRules(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("rules state last sent=%s\n", sent.Format(time.RFC3339))
	diff := rules.Diff(app.Rules)
	log.Println(diff.String())
	return diff, nil
}

func (app *App) Deploy(ctx context.Context, dryRun bool) error {
	diff, err := app.diff(ctx)
	if err != nil {
		return err
	}
	addRules := make(Rules, 0, len(diff))
	deleteIDs := make([]string, 0, len(diff))
	for _, d := range diff {
		if d.Add {
			addRules = append(addRules, Rule{
				Tag:   d.Rule.Tag,
				Value: d.Rule.Value,
			})
			continue
		}
		if d.Delete {
			deleteIDs = append(deleteIDs, d.Rule.ID)
		}
	}
	if err := app.Client.PostRules(ctx, nil, deleteIDs, dryRun); err != nil {
		return err
	}
	if err := app.Client.PostRules(ctx, addRules, nil, dryRun); err != nil {
		return err
	}
	return nil
}
