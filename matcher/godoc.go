package matcher

import (
	"errors"
	"strings"
)

type GodocMatcher struct{}

func NewGodocMatcher() (Matcher, error) {
	return &GodocMatcher{}, nil
}

func (m *GodocMatcher) Match(codes, comments []string) (string, error) {
	matched := []string{}

	in := false
	for _, c := range codes {
		s := strings.Trim(c, " \t")
		if !in && s == "" {
			matched = []string{}
			continue
		}
		if s == "" {
			matched = append(matched, s)
			continue
		}
		if !in && strings.HasPrefix(s, "package ") {
			if len(matched) > 0 {
				return strings.Join(matched, "\n"), nil
			}
			return "", errors.New("no matches")
		}
		if s == "/*" {
			in = true
			matched = []string{}
			continue
		}
		if s == "*/" {
			in = false
			continue
		}

		if strings.HasPrefix(s, "/*") {
			in = true
			matched = append(matched, strings.TrimSuffix(strings.Trim(strings.TrimPrefix(s, "/*"), " \t"), "*/"))
			continue
		}
		if strings.HasSuffix(s, "*/") {
			in = false
			matched = append(matched, strings.Trim(strings.TrimSuffix(s, "/*"), " \t"))
			continue
		}
		if !in {
			if strings.HasPrefix(s, "//") {
				matched = append(matched, strings.Trim(strings.TrimPrefix(s, "//"), " \t"))
			}
			continue
		}
		matched = append(matched, s)
	}

	return "", errors.New("no matches")
}
