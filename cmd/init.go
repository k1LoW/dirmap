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
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/k1LoW/dirmap/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: fmt.Sprintf("Generate dirmap config file (%s)", config.DefaultConfigFilePaths[0]),
	Long:  fmt.Sprintf("Generate dirmap config file (%s).", config.DefaultConfigFilePaths[0]),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		if err := c.Load(""); err != nil {
			return err
		}
		if c.Loaded() {
			return errors.New("dirmap config file is already exist")
		}
		f, err := os.OpenFile(config.DefaultConfigFilePaths[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
		if err != nil {
			return err
		}
		encoder := yaml.NewEncoder(f)
		if err := encoder.Encode(c); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		cmd.Printf("Create %s\n", config.DefaultConfigFilePaths[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
