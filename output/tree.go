package output

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/k1LoW/dirmap/scanner"
	"github.com/xlab/treeprint"
)

func WriteTree(fsys fs.FS, wr io.Writer) error {
	tree := treeprint.New()
	nodes := map[string]treeprint.Tree{}
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
		name := fmt.Sprintf("%s/", d.Name())
		if di.Sys() != nil && !di.Sys().(*scanner.DirInfo).Ignore {
			dirInfo := di.Sys().(*scanner.DirInfo)
			o := strings.Split(dirInfo.Overview, "\n")
			oo := ""
			for _, l := range o {
				if l == "" {
					continue
				}
				if !strings.HasSuffix(l, ".") && !strings.HasSuffix(l, ",") && !strings.HasSuffix(l, "。") && !strings.HasSuffix(l, "、") {
					l = fmt.Sprintf("%s.", l)
				}
				if oo == "" {
					oo = l
				} else {
					oo = fmt.Sprintf("%s %s", oo, l)
				}
			}
			name = fmt.Sprintf("%s/ ... %s", d.Name(), oo)
		}
		p, ok := nodes[filepath.Dir(path)]
		if ok {
			nodes[path] = p.AddBranch(name)
		} else {
			nodes[path] = tree.AddBranch(name)
		}
		return nil
	}); err != nil {
		return err
	}
	if _, err := fmt.Fprint(wr, tree.String()); err != nil {
		return err
	}
	return nil
}
