package twstrulemgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Endpoint    string
	BearerToken string
	RulesFile   string

	Rules Rules
}

func (config *Config) Validate() error {
	if config.BearerToken == "" {
		return errors.New("BearerToken is required")
	}
	if config.Rules == nil {
		if config.RulesFile == "" {
			return errors.New("RulesFile is required")
		}
		if err := config.readRules(); err != nil {
			return err
		}
	}
	if err := config.Rules.Validate(); err != nil {
		return fmt.Errorf("rule file validate failed: %w", err)
	}
	if config.Endpoint == "" {
		config.Endpoint = "api.twitter.com"
	}
	return nil
}

func (config *Config) readRules() error {
	fp, err := os.Open(config.RulesFile)
	if err != nil {
		return fmt.Errorf("can not open rule file: %w", err)
	}
	defer fp.Close()
	var data RulesFile
	decoder := json.NewDecoder(fp)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("can not parse rule file: %w", err)
	}
	config.Rules = data.Rules
	return nil
}
