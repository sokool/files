package files

import (
	"github.com/sokool/domain"
)

type Location = domain.Path

func NewLocation(path string) (Location, error) {
	return domain.NewPath(path)
}

type Meta map[string]string

func (m Meta) merge(n map[string]string) error {
	for key, value := range n {
		if m[key] != "" {
			return Err("%s key already exists", key)
		}
		m[key] = value
	}
	return nil
}
