package files

import (
	"encoding/json"
	"io"

	"github.com/sokool/domain"
)

type Location = domain.Path

func NewLocation(path string) (Location, error) {
	return domain.NewPath(path)
}

type Meta map[string]string

func (m Meta) Merge(n map[string]string) error {
	for key, value := range n {
		if m[key] != "" {
			return Err("%s key already exists", key)
		}
		m[key] = value
	}
	return nil
}

func (m Meta) Map(s Meta) Meta {
	if len(s) == 0 {
		return m
	}
	n := make(Meta)
	for mk := range m {
		for sk, sv := range s {
			if sk == mk {
				if sv == "" {
					sv = mk
				}
				n[sv] = m[mk]
				break
			}
		}
	}
	return n
}

func (m Meta) Size() int {
	return len(m)
}

func (m Meta) WriteTo(w io.Writer) (int64, error) {
	var v any = m
	if m.Size() == 1 {
		for n := range m {
			v = m[n]
			break
		}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}
