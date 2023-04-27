package files

import (
	"io"

	"github.com/sokool/domain"
)

type Storage interface {
	Write(Location, io.Reader) error
	Read(Location, io.Writer) error
	Files(d Location, recursive ...bool) ([]string, error)
}

func NewStorage(url string) (Storage, error) {
	var u domain.URL
	var err error
	if u, err = domain.NewURL(url); err != nil {
		return nil, Err("files: invalid `%s` url", url)
	}
	switch u.Schema {
	case "s3":
		return newS3(u)
	case "memory":
		return make(memory), nil
	default:
		return nil, Err("files: %s schema not supported", u.Schema)
	}
}
