/*
Copyright © 2021 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/k1LoW/dirmap/config"
	"github.com/k1LoW/dirmap/output"
	"github.com/k1LoW/dirmap/scanner"
	"github.com/spf13/cobra"
)

var (
	configPath string
	outPath    string
	format     string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [DIR]",
	Short: "generate directory map",
	Long:  `generate directory map.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		if err := c.Load(configPath); err != nil {
			return err
		}

		dir := ""
		if len(args) == 0 {
			dir = "."
		} else {
			dir = args[0]
		}
		fsys := os.DirFS(dir)
		if err := c.LoadGitIgnore(fsys); err != nil {
			return err
		}
		dmfs, err := scanner.Scan(c, fsys)
		if err != nil {
			return err
		}

		var wr io.Writer
		if outPath != "" {
			file, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				return err
			}
			wr = file
		} else {
			wr = os.Stdout
		}

		switch format {
		case output.Tree:
			if err := output.WriteTree(dmfs, wr); err != nil {
				return err
			}
		case output.Table:
			if err := output.WriteTable(dmfs, wr); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid format: %s", format)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	generateCmd.Flags().StringVarP(&outPath, "out", "o", "", "output file path")
	generateCmd.Flags().StringVarP(&format, "format", "t", output.Tree, "output format")
}
