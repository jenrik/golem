package golem

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
	"net/http"
)

type Context struct {
	Link  string
	JobId string
	Defs  *[]Definition
	HTTP  *http.Response
	HTML  *html.Node
}

// Download and extract data from a page
func ScrapePage(task *Task, defs *[]Definition) (map[string][]map[string]interface{}, error) {
	logCtx := log.WithFields(log.Fields{
		"jobid": task.JobId,
		"link":  task.Link,
	})
	// Find definitions that matches link
	var mdefs []*Definition
	for _, def := range *defs {
		if def.MatchPage(&task.Link) {
			mdefs = append(mdefs, &def)
		}
	}
	if len(mdefs) == 0 {
		logCtx.Error("No definition matched")
		return nil, nil
	}

	// Download page
	resp, err := http.Get(task.Link)
	if err != nil {
		logCtx.WithError(err).Warn("Failed to download page")
		return nil, err
	}
	defer resp.Body.Close()
	// Parse page
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logCtx.WithError(err).Warn("Failed to parse page")
		return nil, err
	}

	// Run all definitions
	var ctx = Context{
		Link:  task.Link,
		JobId: task.JobId,
		Defs:  defs,
		HTTP:  resp,
		HTML:  doc,
	}
	var defData = make(map[string][]map[string]interface{})
	for _, def := range mdefs {
		data, err := runDefinition(def, ctx)
		if err != nil {
			logCtx.WithFields(log.Fields{
				"definition": *def,
			}).WithError(err).Warn("A definition run failed")
			return nil, err
		}
		defData[def.Name] = data
	}

	return defData, nil
}

// Extract data from a page
func runDefinition(def *Definition, ctx Context) ([]map[string]interface{}, error) {
	var data = make([]map[string]interface{}, len(def.Stages))
	for i, stage := range def.Stages {
		data[i] = make(map[string]interface{})
		for _, actionconf := range stage {
			// Get action
			action, ok := actions[actionconf.Action]
			if !ok {
				return nil, errors.New("Encountered a non existing action")
			}

			// Expand config
			config, err := action.ExpandConfig(&data, &ctx, actionconf.Config)
			if err != nil {
				return nil, err
			}

			// Run action
			data[i][actionconf.Name], err = action.Action(&ctx, config)
			if err != nil {
				return nil, err
			}
		}
	}
	return data, nil
}
