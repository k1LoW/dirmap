// Configuration file
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/k1LoW/dirmap/matcher"
	"github.com/k1LoW/expand"
)

var DefaultConfigFilePaths = []string{".dirmap.yml", "dirmap.yml"}

type Config struct {
	Targets Targets  `yaml:"targets"`
	Ignores []string `yaml:"ignores,omitempty"`
	// working directory
	wd string
	// config file path
	path string
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
