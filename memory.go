package files

import (
	"bytes"
	"io"
)

type memory map[Location]*bytes.Buffer

func (m memory) Read(n Location, to io.Writer) error {
	g, ok := m[n]
	if !ok {
		return Err("file %s not found", n)
	}
	_, err := io.Copy(to, g)
	return err
}

func (m memory) Write(n Location, from io.Reader) error {
	var b bytes.Buffer
	if _, err := io.Copy(&b, from); err != nil {
		return err
	}
	m[n] = &b
	return nil
}

func (m memory) Files(d Location, recursive ...bool) ([]string, error) {
	var s []string
	for n := range m {
		s = append(s, n.String())
	}
	return s, nil
}
