package files

import (
	"fmt"
	"io"

	"github.com/sokool/domain"
)

type Service struct{ s Storage }

func NewServiceFromStorage(s Storage) *Service {
	return &Service{s}
}

func NewService(url string) (*Service, error) {
	e, err := NewStorage(url)
	if err != nil {
		return nil, err
	}
	return NewServiceFromStorage(e), nil
}

func (s *Service) Uploader(filepath string) *Uploader {
	return NewUploader(s.s).Filename(filepath)
}

func (s *Service) Downloader(filepath string) (*Downloader, error) {
	f, err := NewLocation(filepath)
	if err != nil {
		return nil, err
	}
	return NewDownloader(s.s, f), nil
}

func (s *Service) Read(filepath string, to io.Writer, m ...Meta) error {
	f, err := NewLocation(filepath)
	if err != nil {
		return err
	}
	return s.s.Read(f, to, m...)
}

func (s *Service) Write(filepath string, from io.Reader, m ...Meta) error {
	f, err := NewLocation(filepath)
	if err != nil {
		return err
	}
	return s.s.Write(f, from, m...)
}

func (s *Service) Files(dir string) ([]string, error) {
	d, err := NewLocation(dir)
	if err != nil {
		return nil, Err("%w", err)
	}
	return s.s.Files(d, true)
}

func (s *Service) String() string {
	var t string
	ff, _ := s.Files("/")
	for i := range ff {
		t += fmt.Sprintf("%s\n", ff[i])
	}
	return t
}

var Err = domain.Errorf
