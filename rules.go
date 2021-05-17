package twstrulemgr

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Rule struct {
	ID    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

func (rule Rule) String() string {
	bs, err := json.Marshal(rule)
	if err != nil {
		return fmt.Sprintf("{ err=%s }", err)
	}
	return string(bs)
}

func (rule Rule) Conflict(o Rule) bool {
	if rule.Value == o.Value {
		return true
	}
	if rule.Tag == "" && o.Tag == "" {
		return false
	}
	if rule.Tag != o.Tag {
		return false
	}
	return rule.Value != o.Value
}

func (rule Rule) Equal(o Rule) bool {
	return rule.Value == o.Value && rule.Tag == o.Tag
}

type Rules []Rule

func (rules Rules) Clone() Rules {
	ret := make(Rules, len(rules))
	copy(ret, rules)
	return ret
}

func (rules Rules) Validate() error {
	for i, l := range rules {
		for j := i + 1; j < len(rules); j++ {
			r := rules[j]
			if l.Value == "" {
				return fmt.Errorf("%s value is required", l)
			}
			if l.Conflict(r) {
				return fmt.Errorf("%s and %s conflict", l, r)
			}
		}
	}
	return nil
}

func (rules Rules) Diff(other Rules) DiffRules {
	ls := rules.Clone()
	rs := other.Clone()
	diff := make(DiffRules, 0, len(rules))
	for i := 0; i < 2; i++ {
		for _, l := range ls {
			found := false
			replace := false
			for _, r := range rs {
				if !l.Conflict(r) {
					continue
				}
				replace = !l.Equal(r)
				found = true
				break
			}
			if !found || replace {
				dr := DiffRule{
					Rule: l,
				}
				if i == 0 {
					dr.Delete = true
				} else {
					dr.Add = true
				}
				diff = append(diff, dr)
			} else if i == 0 {
				dr := DiffRule{
					Rule: l,
				}
				diff = append(diff, dr)
			}

		}
		if i == 0 {
			ls, rs = rs, ls
		}
	}
	sort.Slice(diff, func(i, j int) bool {
		iValue := diff[i].Rule.Value
		jValue := diff[j].Rule.Value
		if strings.HasPrefix(iValue, jValue) || strings.HasPrefix(jValue, iValue) {
			//先頭文字がいくつか一致しているなら、 Delete => Addの順にしたい
			if diff[i].Delete && !diff[j].Delete {
				return true
			}
			return false
		}
		//それ以外は、Valueの順序で
		if diff[i].Rule.Value != diff[j].Rule.Value {
			return diff[i].Rule.Value < diff[j].Rule.Value
		}

		return false
	})
	return diff
}

type DiffRule struct {
	Add    bool
	Delete bool
	Rule   Rule
}

func (diff DiffRule) String() string {
	var str string
	switch {
	case diff.Add:
		str = "+"
	case diff.Delete:
		str = "-"
	default:
		str = " "
	}
	return fmt.Sprintf("%s %s", str, diff.Rule)
}

type DiffRules []DiffRule

func (diffs DiffRules) String() string {
	str := "["
	count := 0
	for _, diff := range diffs {
		str += "\n  " + diff.String()
		count++
	}
	if count > 0 {
		str += "\n"
	}
	str += "]"
	return str
}

type RulesFile struct {
	Rules Rules `json:"rules"`
}
