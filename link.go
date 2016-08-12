package golem

import (
	"net/url"
	"strings"
)

type ActionLink struct{}

func (action *ActionLink) Action(ctx *Context, config interface{}) (interface{}, error) {
	if linkz, ok := config.([]string); ok {
		var links []string
		rel, _ := url.Parse(ctx.Link)
		for _, l := range linkz {
			if link, ok := action.cleanLink(&l, rel, ctx); ok {
				links = append(links, link)
			}
		}
		if len(links) > 0 {
			SubmitTask(ctx.JobId, &links)
		}
	} else if link, ok := config.(string); ok {
		rel, _ := url.Parse(ctx.Link)
		if l, ok := action.cleanLink(&link, rel, ctx); ok {
			SubmitTask(ctx.JobId, &[]string{l})
		}
	}

	return nil, nil
}

func (action *ActionLink) cleanLink(link *string, rel *url.URL, ctx *Context) (string, bool) {
	if l, err := url.Parse(*link); err == nil {
		resolved := ResolveLink(l, rel)
		lc := resolved.String()
		for _, def := range *ctx.Defs {
			if def.MatchPage(&lc) {
				return lc, true
			}
		}
	}
	return "", false
}

func ResolveLink(link *url.URL, rel *url.URL) url.URL {
	if len(link.Host) == 0 {
		link.Host = rel.Host
	}
	if len(link.Scheme) == 0 {
		link.Scheme = rel.Scheme
	}
	if rel.User != nil && link.User == nil {
		link.User = rel.User
	}
	// Resolve relative path
	if len(link.Path) > 0 {
		path := link.EscapedPath()
		relPath := rel.EscapedPath()
		if len(path) >= 3 && path[0:3] == "../" { // Down
			var counter int
			for len(path) >= 3 && path[0:3] == "../" {
				counter += 1
				path = path[3:]
			}
			segments := strings.Split(relPath, "/")
			var lastSlash bool = false
			if len(segments) > 0 && segments[len(segments)-1] == "" {
				segments = segments[:len(segments)-2]
				lastSlash = true
			}
			if counter >= len(segments) {
				if lastSlash {
					path = "/" + path
				}
				p, _ := url.QueryUnescape(path)
				link.Path = p
			} else {
				if lastSlash {
					path = "/" + path
				}
				pos := len(segments) - counter
				segments = segments[:pos]
				if len(path) > 0 && path[0] != 47 {
					path = "/" + path
				}
				p, _ := url.QueryUnescape(strings.Join(segments, "/") + path)
				link.Path = p
				if len(link.Path) > 0 && link.Path[0] != 47 {
					link.Path = "/" + link.Path
				}
			}
		} else { // Up and absolute
			if path[0] != 47 { // Up
				if relPath[len(relPath)-1] != 47 {
					path = relPath + "/" + path
				} else {
					path = relPath + path
				}
			} // Else absolute
			link.Path, _ = url.QueryUnescape(path)
		}
	}

	return *link
}

func (action *ActionLink) ExpandConfig(stages *[]map[string]interface{}, ctx *Context, config interface{}) (interface{}, error) {
	conf := config.(map[string]interface{})
	template := conf["link"].(string)
	return TemplateValue(stages, ctx, &template)
}

func (action *ActionLink) Validate(config interface{}) bool {
	if config, ok := config.(map[string]interface{}); ok {
		_, ok := config["link"]
		return ok
	} else {
		return false
	}
}

func init() {
	RegisterAction(new(ActionLink), "link")
}
