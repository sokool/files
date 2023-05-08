package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/sokool/domain"
)

type Uploader struct {
	filespath string
	storage   Storage
	extension bool
	size      int64
	key       string
	media     map[string]bool
	request   Meta
	response  func(Location, Meta) any
}

func NewUploader(s Storage) *Uploader {
	return &Uploader{
		storage: s,
		media:   make(map[string]bool),
		key:     "file",
	}
}

func (u *Uploader) Form(key string) *Uploader {
	u.key = key
	return u
}

func (u *Uploader) MaxSize(n int64) *Uploader {
	u.size = n
	return u
}

func (u *Uploader) AllowedMedia(types ...string) *Uploader {
	for i := range types {
		u.media[types[i]] = true
	}
	return u
}

func (u *Uploader) FilesPath(s string) *Uploader {
	u.filespath = s
	return u
}

func (u *Uploader) Request(m Meta) *Uploader {
	u.request = m
	return u
}

func (u *Uploader) Response(f func(l Location, m Meta) any) *Uploader {
	u.response = f
	return u
}

func (u *Uploader) AutoExtension(b bool) *Uploader {
	u.extension = b

	//detectedFileType := http.DetectContentType(fileBytes)
	//switch detectedFileType {
	//case "image/jpeg", "image/png":
	//	break
	//default:
	//	log.WithError(err).Debug("invalid file type")
	//	render.Render(w, r, contract.ErrInvalidRequest(err))
	//	return
	//}
	//
	//fileEndings, err := mime.ExtensionsByType(detectedFileType)
	//if err != nil {
	//	log.WithError(err).Error("failed to detect file extension for specified file type")
	//	render.Render(w, r, contract.ErrInternalServer)
	//	return
	//}
	return u
}

func (u *Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh, h, err := r.FormFile(u.key)
	if err != nil {
		err = Err("src multipart file not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer fh.Close()

	if len(u.media) != 0 {
		t := h.Header.Get("content-type")
		if _, ok := u.media[t]; !ok && len(u.media) != 0 {
			err = Err("%s content type not supported", t)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if u.size > 0 && h.Size > u.size {
		err = Err("%d file size too big, max %d", h.Size, u.size)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l, f, err := u.file(fh)
	if err != nil {
		http.Error(w, Err("invalid file name %w", err).Error(), http.StatusBadRequest)
		return
	}

	m := make(Meta)
	if u.request != nil {
		m = u.request
	}
	if err = u.storage.Write(l, f, m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.response != nil {
		if err = json.NewEncoder(w).Encode(u.response(l, m)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (u *Uploader) file(r io.Reader) (Location, io.Reader, error) {
	var err error
	var l Location
	var d domain.ID
	var f string
	if err = domain.NewID(&d); err != nil {
		return l, r, err
	}
	if f = u.filespath + "/" + d.String(); u.extension {
		var b []byte
		var e []string
		if b, err = io.ReadAll(r); err != nil {
			return l, r, err
		}

		if e, err = mime.ExtensionsByType(http.DetectContentType(b)); err != nil {
			return l, r, err
		}
		if len(e) == 1 {
			f = fmt.Sprintf("%s%s", f, e[0])
		}
		r = bytes.NewBuffer(b)
	}

	l, err = NewLocation(f)
	return l, r, err
}
