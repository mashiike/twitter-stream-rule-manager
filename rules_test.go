package twstrulemgr_test

import (
	"testing"

	twstrulemgr "github.com/mashiike/twitter-stream-rule-manager"
	"github.com/stretchr/testify/assert"
)

func TestRulesDiff(t *testing.T) {
	leftRules := twstrulemgr.Rules{
		{
			// stable
			ID:    "1234567890",
			Tag:   "hoge",
			Value: "1-hoge",
		},
		{
			// stable value only
			ID:    "2345678901",
			Value: "2-fuga",
		},
		{
			// replace tag add
			ID:    "3456789012",
			Value: "3-piyo",
		},
		{
			// replace same tag , value change
			ID:    "4567890123",
			Value: "4-tora",
			Tag:   "tora",
		},
		{
			// replace same value , tag change
			ID:    "5678901234",
			Value: "5-tama",
			Tag:   "tama",
		},
		{
			// delete rule
			ID:    "6789012345",
			Value: "6-hage",
			Tag:   "hage",
		},
	}
	rightRules := twstrulemgr.Rules{
		{
			Tag:   "hoge",
			Value: "1-hoge",
		},
		{
			Value: "2-fuga",
		},
		{
			Tag:   "piyo",
			Value: "3-piyo",
		},
		{
			Tag:   "tora",
			Value: "4-tora -RT",
		},
		{
			Tag:   "tama nyan",
			Value: "5-tama",
		},
		{
			// add new rule
			Tag:   "poyo",
			Value: "7-poyo",
		},
	}
	expected := `[
    {"id":"1234567890","value":"1-hoge","tag":"hoge"}
    {"id":"2345678901","value":"2-fuga"}
  - {"id":"3456789012","value":"3-piyo"}
  + {"value":"3-piyo","tag":"piyo"}
  - {"id":"4567890123","value":"4-tora","tag":"tora"}
  + {"value":"4-tora -RT","tag":"tora"}
  - {"id":"5678901234","value":"5-tama","tag":"tama"}
  + {"value":"5-tama","tag":"tama nyan"}
  - {"id":"6789012345","value":"6-hage","tag":"hage"}
  + {"value":"7-poyo","tag":"poyo"}
]`
	assert.Equal(t, expected, leftRules.Diff(rightRules).String())
}
