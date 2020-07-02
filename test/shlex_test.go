// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shlex_test

import (
	"fmt"
	"reflect"
	"testing"

	anmitsu "github.com/anmitsu/go-shlex"
	goog "github.com/google/shlex"
	"github.com/hugelgupf/go-shlex"
)

func TestArgv(t *testing.T) {
	for i, tt := range []struct {
		desc         string
		in           string
		want         []string
		anmitsuWrong bool
		googWrong    bool
	}{
		{
			desc: "nothing to do",
			in:   "stuff",
			want: []string{"stuff"},
		},
		{
			desc: "split",
			in:   "stuff more stuff",
			want: []string{"stuff", "more", "stuff"},
		},
		{
			desc: "escape",
			in:   "stuff\\ more stuff",
			want: []string{"stuff more", "stuff"},
		},
		{
			desc: "quote",
			in:   "stuff var='more stuff'",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "double quote",
			in:   "stuff var=\"more stuff\"",
			want: []string{"stuff", "var=more stuff"},
		},
		{
			desc: "single quote specials",
			in:   `stuff var='more stuff $ \ \$ " ` + "` '",
			want: []string{"stuff", `var=more stuff $ \ \$ " ` + "` "},
		},
		{
			desc:      "quote forgot close",
			in:        "stuff var='more stuff",
			want:      []string{"stuff", "var=more stuff"},
			googWrong: true,
		},
		{
			desc: "empty",
			in:   "",
			want: []string{},
		},

		// GNU Bash manual:
		//
		// The backslash retains its special meaning only when followed
		// by one of the following characters: ‘$’, ‘`’, ‘"’, ‘\’, or
		// newline. Within double quotes, backslashes that are followed
		// by one of these characters are removed.
		{
			desc:      "double quote newline",
			in:        `stuff var="more stuff \n "`,
			want:      []string{"stuff", `var=more stuff \n `},
			googWrong: true,
		},
		{
			desc:      "double quote backslash",
			in:        `stuff var="more stuff \ "`,
			want:      []string{"stuff", `var=more stuff \ `},
			googWrong: true,
		},
		{
			desc: "double quote double quote",
			in:   `stuff var="more stuff \" \" \\\" "`,
			want: []string{"stuff", `var=more stuff " " \" `},
		},
		{
			desc:         "double quote dollar",
			in:           `stuff var="more stuff $ \$ \\$ \\\$ "`,
			want:         []string{"stuff", `var=more stuff $ $ \$ \$ `},
			anmitsuWrong: true,
		},
		{
			desc:         "double quote special backtick",
			in:           "stuff var=\"more stuff ` \\`\"",
			want:         []string{"stuff", "var=more stuff ` `"},
			anmitsuWrong: true,
		},

		// (Fixed) tests from anmitsu/go-shlex
		{
			in: `This string has an embedded apostrophe, doesn't it?`,
			want: []string{
				"This",
				"string",
				"has",
				"an",
				"embedded",
				"apostrophe,",
				"doesnt it?",
			},
			googWrong: true,
		},
		{
			in: "This string has embedded \"double quotes\" and 'single quotes' in it,\nand even \"a 'nested example'\".\n",
			want: []string{
				"This",
				"string",
				"has",
				"embedded",
				`double quotes`,
				"and",
				`single quotes`,
				"in",
				"it,",
				"and",
				"even",
				`a 'nested example'.`,
			},
		},
		{
			in: `Hello world!, こんにちは　世界！`,
			want: []string{
				"Hello",
				"world!,",
				"こんにちは",
				"世界！",
			},
			googWrong: true,
		},
		{
			in:   `Do"Not"Separate`,
			want: []string{`DoNotSeparate`},
		},
		{
			in: `Escaped \e Character not in quotes`,
			want: []string{
				"Escaped",
				"e",
				"Character",
				"not",
				"in",
				"quotes",
			},
		},
		{
			in: `Escaped "\e" Character in double quotes`,
			want: []string{
				"Escaped",
				`\e`,
				"Character",
				"in",
				"double",
				"quotes",
			},
			googWrong: true,
		},
		{
			in: `Escaped '\e' Character in single quotes`,
			want: []string{
				"Escaped",
				`\e`,
				"Character",
				"in",
				"single",
				"quotes",
			},
		},
		{
			in: `Escaped '\'' \"\'\" single quote`,
			want: []string{
				"Escaped",
				`\ \"\"`,
				"single",
				"quote",
			},
		},
		{
			in: `Escaped "\"" \'\"\' double quote`,
			want: []string{
				"Escaped",
				`"`,
				`'"'`,
				"double",
				"quote",
			},
		},
		{
			in:   `"'Strip extra layer of quotes'"`,
			want: []string{`'Strip extra layer of quotes'`},
		},
		// Tests from google/shlex
		{
			in:           "one two \"three four\" \"five \\\"six\\\"\" seven#eight # nine # ten\n eleven 'twelve\\' thirteen=13 fourteen/14",
			want:         []string{"one", "two", "three four", "five \"six\"", "seven#eight", "eleven", "twelve\\", "thirteen=13", "fourteen/14"},
			anmitsuWrong: true,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			got := shlex.Split(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split = %#v, want %#v", got, tt.want)
			}

			got, _ = goog.Split(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				if tt.googWrong {
					t.Logf("(wrong) google/shlex Split = %#v, want %#v", got, tt.want)
				} else {
					t.Errorf("google/shlex Split = %#v, want %#v", got, tt.want)
				}
			} else if tt.googWrong {
				t.Errorf("google/shlex is right, but marked wrong")
			}

			got, _ = anmitsu.Split(tt.in, true)
			if !reflect.DeepEqual(got, tt.want) {
				if tt.anmitsuWrong {
					t.Logf("(wrong) anmitsu/go-shlex Split = %#v, want %#v", got, tt.want)
				} else {
					t.Errorf("anmitsu/go-shlex Split = %#v, want %#v", got, tt.want)
				}
			} else if tt.anmitsuWrong {
				t.Errorf("anmitsu/go-shlex is right, but marked wrong")
			}
		})
	}
}
