package matcher

import (
	"errors"
	"regexp"
	"strings"
)

type RegexpMatcher struct {
	re *regexp.Regexp
}

func NewRegexpMatcher(expr string) (Matcher, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &RegexpMatcher{
		re: re,
	}, nil
}

func (m *RegexpMatcher) Match(codes, comments []string) (string, error) {
	for _, c := range codes {
		s := strings.Trim(c, " \t")
		if s == "" {
			continue
		}
		matches := m.re.FindAllStringSubmatch(s, -1)
		if len(matches) == 0 {
			continue
		}
		if len(matches[0]) > 1 {
			return matches[0][1], nil
		} else {
			return matches[0][0], nil
		}
	}
	return "", errors.New("no matches")
}
