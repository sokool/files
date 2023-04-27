package files

import (
	"regexp"
	"strings"
)

type Location struct{ string }

func NewLocation(path string) (Location, error) {
	if path == "/" {
		return Location{path}, nil
	}
	if ok, _ := regexp.MatchString(`^(/[\w\s~+!'.-]+)+$`, path); ok {
		return Location{path}, nil
	}
	return Location{}, Err("`%s` must start from / character followed by alphanumerics and/or _ ' ! . - ", path)
}

// Cut first n character from underlying location string
func (l Location) Cut(n int) string {
	if l.IsZero() {
		return ""
	}
	return l.string[n:]
}

// Tail gives last part of location, it might be directory or file with extension
func (l Location) Tail() string {
	s := strings.Split(l.string, "/")
	return s[len(s)-1]
}

func (l Location) IsZero() bool {
	return l.string == ""
}

func (l Location) String() string {
	return l.string
}
