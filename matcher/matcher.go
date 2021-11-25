package matcher

const (
	Markdown        = "markdown"
	MarkdownHeading = "markdownHeading"
	Godoc           = "godoc"
)

var MatcherMap = map[string]func() (Matcher, error){
	"markdown":        NewMarkdownMatcher,
	"markdownHeading": NewMarkdownHeadingMatcher,
	"godoc":           NewGodocMatcher,
}

type Matcher interface {
	Match(codes, comments []string) (string, error)
}
