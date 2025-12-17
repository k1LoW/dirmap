package output

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/k1LoW/dirmap/scanner"
)

func WriteTable(fsys fs.FS, wr io.Writer) error {
	if _, err := fmt.Fprintf(wr, "| %s |\n", strings.Join([]string{"Directory", "Overview"}, " | ")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(wr, "| %s |\n", strings.Join([]string{"---", "---"}, " | ")); err != nil {
		return err
	}
	if err := fs.WalkDir(fsys, scanner.RootKey, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == scanner.RootKey {
			return nil
		}
		di, err := d.Info()
		if err != nil {
			return err
		}
		name := strings.TrimPrefix(strings.Replace(fmt.Sprintf("%s/", path), scanner.RootKey, "", 1), "/")
		o := ""
		if di.Sys() != nil {
			dirInfo, ok := di.Sys().(*scanner.DirInfo)
			if !ok {
				return nil
			}
			if dirInfo.Ignore {
				return nil
			}
			o = fmt.Sprintf("%s ( [ref](%s) )", strings.ReplaceAll(dirInfo.Overview, "\n", "<br>"), dirInfo.File)
		}
		if _, err := fmt.Fprintf(wr, "| %s |\n", strings.Join([]string{name, o}, " | ")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
