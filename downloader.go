package files

import "net/http"

type Downloader struct {
	file    Location
	storage Storage
}

func NewDownloader(s Storage, l Location) *Downloader {
	return &Downloader{storage: s, file: l}
}

func (d *Downloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := d.storage.Read(d.file, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
