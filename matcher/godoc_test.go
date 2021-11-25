package matcher

import (
	"strings"
	"testing"
)

func TestGodocMatcher(t *testing.T) {
	tests := []struct {
		code    string
		want    string
		wantErr bool
	}{
		{"", "", true},
		{
			`/*
dirmap
This is overview

Hello dirmap
*/
package cmd
`,
			`dirmap
This is overview

Hello dirmap`,
			false,
		},
		{
			`
package cmd

/*
dirmap
This is overview

Hello dirmap
*/
`,
			"",
			true,
		},
		{
			`
// Hello dirmap

// dirmap
// This is overview
package cmd
`,
			`dirmap
This is overview`,
			false,
		},
		{
			`
// Hello dirmap
// 
// dirmap
// This is overview
package cmd
`,
			`Hello dirmap

dirmap
This is overview`,
			false,
		},
	}
	for _, tt := range tests {
		m, _ := NewGodocMatcher()
		codes := strings.Split(tt.code, "\n")
		got, err := m.Match(codes, []string{})
		if err != nil {
			if !tt.wantErr {
				t.Error(err)
			}
			continue
		}
		if tt.wantErr {
			t.Error("should be error")
		}
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
