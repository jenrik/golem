package golem

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

func TemplateInsert(stage *[]map[string]interface{}, ctx *Context, template *string) (interface{}, error) {
	// ToDo or not
	return nil, nil
}

func TemplateValue(stages *[]map[string]interface{}, ctx *Context, template *string) (interface{}, error) {
	segments := strings.Split(*template, "|")
	switch segments[0] {
	case "stage":
		if len(segments) > 3 || len(segments) < 2 {
			err := errors.New("Received invalid amount of arguments to stage in template")
			log.WithFields(log.Fields{
				"template": *template,
				"amount":   len(segments),
			}).Warn(err.Error())
			return nil, err
		}
		if i, err := strconv.Atoi(segments[1]); err == nil {
			if i > 0 || i < len(*stages) {
				stage := (*stages)[i]
				if len(segments) == 2 {
					return stage, nil
				} else {
					if value, ok := stage[segments[2]]; ok {
						return value, nil
					} else {
						err := errors.New("Received non-existing key for data in stage in template")
						log.WithFields(log.Fields{
							"template": *template,
							"key":      segments[2],
						}).Warn(err.Error())
						return nil, err
					}
				}
			} else {
				err := errors.New("Received stage index out of bounds in a template")
				log.WithFields(log.Fields{
					"template": *template,
					"index":    i,
				}).Warn(err.Error())
				return nil, err
			}
		} else {
			err := errors.New("Received non-number as stage index in a template")
			log.WithFields(log.Fields{
				"template": *template,
				"value":    segments[1],
			}).Warn(err.Error())
			return nil, err
		}
	// ToDo case "context":
	default:
		return nil, errors.New("Template requested unknown")
	}
}
