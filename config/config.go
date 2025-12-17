// Configuration file
package config

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/k1LoW/dirmap/matcher"
	"github.com/k1LoW/expand"
	gitignore "github.com/sabhiram/go-gitignore"
)

var DefaultConfigFilePaths = []string{".dirmap.yml", "dirmap.yml"}

type Config struct {
	Targets Targets  `yaml:"targets"`
	Ignores []string `yaml:"ignores,omitempty"`
	// working directory
	wd string
	// config file path
	path string

	GitIgnore *gitignore.GitIgnore `yaml:"-"`
}

type Target struct {
	File          string `yaml:"file"`
	MatcherString string `yaml:"matcher"`
}

func (t *Target) Matcher() (matcher.Matcher, error) {
	mfunc, ok := matcher.MatcherMap[t.MatcherString]
	if ok {
		return mfunc()
	}
	return matcher.NewRegexpMatcher(t.MatcherString)
}

type Targets []*Target

var defaultTargets = Targets{
	&Target{
		File:          ".dirmap.md",
		MatcherString: matcher.Markdown,
	},
	&Target{
		File:          "README.md",
		MatcherString: matcher.Markdown,
	},
	&Target{
		File:          "doc.go",
		MatcherString: matcher.Godoc,
	},
	&Target{
		File:          "*.go",
		MatcherString: matcher.Godoc,
	},
}

func New() *Config {
	wd, _ := os.Getwd()
	return &Config{
		wd: wd,
	}
}

func (c *Config) Getwd() string {
	return c.wd
}

func (c *Config) Setwd(path string) {
	c.wd = path
}

func (c *Config) Load(path string) error {
	if path == "" {
		for _, p := range DefaultConfigFilePaths {
			if f, err := os.Stat(filepath.Join(c.wd, p)); err == nil && !f.IsDir() {
				if path != "" {
					return fmt.Errorf("duplicate config file [%s, %s]", path, p)
				}
				path = p
			}
		}
	}
	if path == "" {
		c.Targets = defaultTargets
		return nil
	}
	c.path = filepath.Join(c.wd, path)
	buf, err := os.ReadFile(filepath.Clean(c.path))
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(expand.ExpandenvYAMLBytes(buf), c); err != nil {
		return err
	}
	return nil
}

func (c *Config) Loaded() bool {
	return c.path != ""
}

// LoadGitIgnore loads all .gitignore files recursively.
func (c *Config) LoadGitIgnore(fsys fs.FS) error {
	// Create a temporary file to store the contents with prefixed path
	tempFile, err := os.CreateTemp("", "dirmap")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	hasGitignore := false

	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) != ".gitignore" {
			return nil
		}
		hasGitignore = true

		// Get the prefix for the gitignore lines
		// If the path is "a/b/.gitignore", the prefix should be "a/b/"
		// Or if the path is just ".gitignore", no prefix is needed
		prefixForGitignoreLines := ""
		lastIndex := strings.LastIndex(path, "/")
		if lastIndex != -1 {
			prefixForGitignoreLines = path[:lastIndex+1]
		}

		// Open the base gitignore file
		baseGitignore, err := os.Open(path)
		if err != nil {
			return err
		}
		defer baseGitignore.Close()

		// Read the base gitignore file and write the contents to the temporary file
		writer := bufio.NewWriter(tempFile)
		scanner := bufio.NewScanner(baseGitignore)

		// We want to skip empty lines, so use scanner.Scan()
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				_, err = writer.WriteString(prefixForGitignoreLines + line + "\n")
				if err != nil {
					return err
				}
			}
		}

		err = writer.Flush()
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if !hasGitignore {
		return nil
	}

	ignore, err := gitignore.CompileIgnoreFile(tempFile.Name())
	if err != nil {
		return err
	}
	c.GitIgnore = ignore
	return nil
}
