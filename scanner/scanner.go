package scanner

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing/fstest"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hhatto/gocloc"
	"github.com/k1LoW/dirmap/config"
)

const RootKey = "root"

var (
	defined = gocloc.NewDefinedLanguages()
)

type DirInfo struct {
	File     string
	Overview string
}

func Scan(c *config.Config, fsys fs.FS) (fstest.MapFS, error) {
	dm := fstest.MapFS{}
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if filepath.Base(path) == ".git" {
			return fs.SkipDir
		}
		for _, p := range c.Ignores {
			match, _ := doublestar.PathMatch(p, path)
			if match {
				return nil
			}
		}
		key := strings.TrimSuffix(fmt.Sprintf("%s/%s", RootKey, strings.TrimSuffix(path, ".")), "/")

		dm[key] = &fstest.MapFile{
			Mode: fs.ModeDir,
			Sys:  nil,
		}
		files, err := fs.ReadDir(fsys, path)
		if err != nil {
			return err
		}
		for _, t := range c.Targets {
			m, err := t.Matcher()
			if err != nil {
				return err
			}
			for _, fi := range files {
				if fi.IsDir() {
					continue
				}
				match, _ := filepath.Match(t.File, fi.Name())
				if !match {
					continue
				}

				p := filepath.Join(path, fi.Name())
				f, err := fsys.Open(p)
				if err != nil {
					return err
				}
				codes, comments := scan(f, p)
				_ = f.Close()
				o, err := m.Match(codes, comments)
				if err != nil {
					continue
				}
				dm[key] = &fstest.MapFile{
					Mode: fs.ModeDir,
					Sys: &DirInfo{
						File:     p,
						Overview: o,
					},
				}
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return dm, nil
}

func scan(r io.Reader, path string) ([]string, []string) {
	opts := gocloc.NewClocOptions()
	codes := []string{}
	comments := []string{}

	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			codes = append(codes, scanner.Text())
		}
		comments = codes
	default:
		l, ok := gocloc.Exts[strings.TrimPrefix(ext, ".")]
		opts.OnBlank = func(line string) {
			codes = append(codes, line)
		}
		opts.OnCode = func(line string) {
			codes = append(codes, line)
		}
		opts.OnComment = func(line string) {
			codes = append(codes, line)
			comments = append(comments, line)
		}
		if ok {
			_ = gocloc.AnalyzeReader("", defined.Langs[l], r, opts)
		} else {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				codes = append(codes, scanner.Text())
			}
			comments = codes
		}
	}
	return codes, comments
}
