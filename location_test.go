package files_test

import (
	"fmt"
	"testing"

	"github.com/sokool/files"
)

func TestNewLocation(t *testing.T) {
	type scenario struct {
		description string
		path        string
		err         bool
	}

	cases := []scenario{
		{"root path is required", "/", false},
		{"strings are fine", "/some/cool/files/path/readme", false},
		{"uppercase are fine", "/HI/THERE", false},
		{"numbers are ok", "/users/35682/file", false},
		{"whitespaces are fine", "/Nice Users/Tom Hilldinor", false},
		{"hyphens are ok", "/documents/invoices-2022-04/fv-1", false},
		{"dashes are ok", "/hello_word/_example_/__fold er__", false},
		{"tilda are ok", "/~Hello_word~2~example", false},
		{"dots are ok", "/some.user.path.with/filename.txt", false},
		{"plus are ok", "/in+flames/guitar+samples", false},
		{"exclamation mark are ok", "/!!important!!!/D.O.C.U.M.E.N.T.S", false},
		{"empty string is not ok", "", true},
		{"dollar sign is not ok", "/$HI$/THERE", true},
		{"multiple slashes are not ok", "/path//name", true},
		{"extra characters at end are not ok", "/dir/game.exe?query=true", true},
	}

	l, err := files.NewLocation("/some/cool/path/file.exe")
	fmt.Println(l, err, l.IsZero(), l.String(), l.Cut(1), l.Tail())
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if _, err := files.NewLocation(c.path); (err != nil && !c.err) || (err == nil && c.err) {
				t.Fatalf("expected error:%v got:%v", c.err, err)
			}
		})
	}
}
