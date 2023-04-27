package files

import (
	"io"
	"net/http"

	"github.com/sokool/domain"
)

type Uploader struct {
	filepath  string
	storage   Storage
	extension bool
	size      int64
	key       string
	media     map[string]bool
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

func (u *Uploader) Filename(filepath string) *Uploader {
	u.filepath = filepath
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
	f, h, err := r.FormFile(u.key)
	if err != nil {
		err = Err("src multipart file not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

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

	l, err := u.location(f)
	if err != nil {
		http.Error(w, Err("invalid file name %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err = u.storage.Write(l, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *Uploader) location(r io.Reader) (Location, error) {
	var err error
	var l Location
	if u.filepath == "" {
		var d domain.ID
		if err = domain.NewID(&d); err != nil {
			return l, err
		}
		u.filepath = d.String()
	}

	if u.extension {
		//mime.ExtensionsByType(http.DetectContentType(fileBytes))
	}
	return NewLocation(u.filepath)
}
