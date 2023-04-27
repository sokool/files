package files

import (
	"net/http"
)

type Uploader struct {
	file    Location
	storage Storage
	size    int64
	media   map[string]bool
}

func NewUploader(s Storage, l Location, size int64, media ...string) *Uploader {
	var u = Uploader{storage: s, file: l, size: size, media: make(map[string]bool, len(media))}
	for i := range media {
		u.media[media[i]] = true
	}
	return &u
}

func (u *Uploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("src")
	if err != nil {
		err = Err("src multipart file not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	t := h.Header.Get("content-type")
	if _, ok := u.media[t]; !ok && len(u.media) != 0 {
		err = Err("%s content type not supported", t)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.Size > u.size {
		err = Err("%d file size too big, max %d", h.Size, u.size)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = u.storage.Write(u.file, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
