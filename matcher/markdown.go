package matcher

import (
	"errors"
	"regexp"
	"strings"
)

type MarkdownMatcher struct{}

func NewMarkdownMatcher() (Matcher, error) {
	return &MarkdownMatcher{}, nil
}

var markdownRe = regexp.MustCompile(`^[#|<!]`)

func (m *MarkdownMatcher) Match(codes, comments []string) (string, error) {
	matched := []string{}

	for _, c := range codes {
		s := strings.Trim(c, " \t")
		if s == "" {
			if len(matched) > 0 {
				return strings.Join(matched, "\n"), nil
			}
			continue
		}
		if markdownRe.MatchString(s) {
			if len(matched) > 0 {
				return strings.Join(matched, "\n"), nil
			}
			continue
		}
		matched = append(matched, s)
	}

	if len(matched) > 0 {
		return strings.Join(matched, "\n"), nil
	}
	return "", errors.New("no matches")
}

func NewMarkdownHeadingMatcher() (Matcher, error) {
	expr := `^#+\s+(.+)$`
	return NewRegexpMatcher(expr)
}
