package golem

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"net/url"
)

type Definition struct {
	Matchers []PageMatcher
	Stages   [][]ActionDef
	Name     string
}

type ActionDef struct {
	Action string
	Name   string
	Config interface{}
}

func ParseDefs(jsonDef *[]byte) (*[]Definition, error) {
	var def []Definition = make([]Definition, 10)
	err := json.Unmarshal(*jsonDef, &def)
	if err != nil {
		log.WithError(err).Warn("Failed to parse definition")
		return nil, err
	} else {
		return &def, nil
	}
}

func (def *Definition) MatchPage(link *string) bool {
	l, err := url.Parse(*link)
	if err != nil {
		log.WithError(err).Error("Failed to parse link before page matching")
		return false
	}
	for _, matcher := range def.Matchers {
		match, err := matcher.Check(l)
		if err != nil {
			log.WithError(err).Error("Error during page matching")
			continue
		}
		if match {
			return true
		}
	}
	return false
}

// Check whether or not the data contained in the definition is valid
func (def *Definition) Validate() (bool, error) {
	for _, stage := range def.Stages {
		for _, actiondef := range stage {
			ok, err := actiondef.Validate()
			if err != nil {
				log.WithError(err).Error("Error during definition validation")
				return false, err
			}
			if !ok {
				log.Error("Invalid definition")
				return false, nil
			}
		}
	}
	return true, nil
}

// Check whether or notthe data container in the Action is valid
func (actionDef *ActionDef) Validate() (bool, error) {
	action, contains := actions[actionDef.Action]
	if !contains {
		return false, errors.New("ActionDef contains a non-existing action")
	}

	return action.Validate(actionDef.Config), nil
}
