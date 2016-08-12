package golem

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type PageMatcher struct {
	Scheme   []string
	User     []string
	Domain   []string
	Path     []map[string]string
	Query    []map[string]string
	Fragment []string
}

// Normalizes the fields
// Currently scheme and domain a forced to lower case
func (matcher *PageMatcher) Normalize() {
	for i, v := range matcher.Scheme {
		matcher.Scheme[i] = strings.ToLower(v)
	}
	for i, v := range matcher.Domain {
		matcher.Domain[i] = strings.ToLower(v)
	}
}

// Check whether or not the link matches
func (matcher *PageMatcher) Check(link *url.URL) (bool, error) {
	if link == nil {
		return true, nil
	}
	matcher.Normalize()

	// Schema
	if matcher.Scheme != nil && len(matcher.Scheme) > 0 && !contains(matcher.Scheme, strings.ToLower(link.Scheme)) {
		return false, nil
	}

	// User
	if matcher.User != nil && len(matcher.User) > 0 && !contains(matcher.User, link.User.String()) {
		return false, nil
	}

	// Domain
	if matcher.Domain != nil && len(matcher.Domain) > 0 && !contains(matcher.Domain, link.Host) {
		return false, nil
	}

	var check = false
	// Path
	if matcher.Path != nil && len(matcher.Path) > 0 {
		var segments = strings.Split(link.RawPath, "/")
		for _, pathmatcher := range matcher.Path {
			// Check max
			smax, ok := pathmatcher["max"]
			if ok {
				max, err := strconv.Atoi(smax)
				if err != nil {
					return false, err
				}
				if max >= len(segments) {
					continue
				}
			}
			// Check min
			smin, ok := pathmatcher["max"]
			var min int
			if ok {
				min, err := strconv.Atoi(smin)
				if err != nil {
					return false, err
				}
				if min >= len(segments) {
					continue
				}
			}
			// Check segment
			var innerCheck = true
			for num, segment := range pathmatcher {
				if n, err := strconv.Atoi(num); err == nil {
					if n < min-1 || n > len(segments)-1 {
						return false, errors.New("Path segment identifier out of bounds")
					}
					if segments[n] != segment {
						innerCheck = false
					}
				}
			}
			if !innerCheck {
				check = true
				break
			}
		}
		if !check {
			return false, nil
		}
	}

	// Query
	if matcher.Query != nil && len(matcher.Query) > 0 {
		check = false
		for _, query := range matcher.Query {
			var all = true
			q, err := url.ParseQuery(link.RawQuery)
			if err != nil {
				return false, err
			}
			for k, v := range query {
				if !contains(q[k], v) {
					all = false
					break
				}
			}
			if all {
				check = true
				break
			}
		}
		if !check {
			return false, nil
		}
	}

	// Fragment
	return matcher.Fragment == nil || len(matcher.Fragment) == 0 || contains(matcher.Fragment, link.Fragment), nil
}

func contains(arr []string, contain string) bool {
	for _, v := range arr {
		if v == contain {
			return true
		}
	}
	return false
}
