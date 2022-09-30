package middleware

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"
)

func FSCacheControl(fsys fs.FS, modTime time.Time, duration time.Duration) http.Handler {
	h := http.FileServer(
		&StaticFSWrapper{
			FileSystem:   http.FS(fsys),
			FixedModTime: modTime,
		},
	)
	maxAge := fmt.Sprintf("max-age=%.0f", duration.Seconds())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", maxAge)
		h.ServeHTTP(w, r)
	})
}

// https://github.com/golang/go/issues/44854#issuecomment-808906568

type StaticFSWrapper struct {
	http.FileSystem
	FixedModTime time.Time
}

func (f *StaticFSWrapper) Open(name string) (http.File, error) {
	file, err := f.FileSystem.Open(name)

	return &StaticFileWrapper{File: file, fixedModTime: f.FixedModTime}, err
}

type StaticFileWrapper struct {
	http.File
	fixedModTime time.Time
}

func (f *StaticFileWrapper) Stat() (os.FileInfo, error) {
	fileInfo, err := f.File.Stat()

	return &StaticFileInfoWrapper{FileInfo: fileInfo, fixedModTime: f.fixedModTime}, err
}

type StaticFileInfoWrapper struct {
	os.FileInfo
	fixedModTime time.Time
}

func (f *StaticFileInfoWrapper) ModTime() time.Time {
	return f.fixedModTime
}
