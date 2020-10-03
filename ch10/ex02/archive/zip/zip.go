package zip

import (
	z "archive/zip"

	"local.gopl/archive"
)

func List(name string) (files []archive.FileInfo, err error) {
	r, err := z.OpenReader(name)
	if err != nil {
		return
	}
	defer r.Close()

	for _, f := range r.File {
		files = append(files, archive.FileInfo{f.FileInfo(), f.Name})
	}
	return
}

func init() {
	archive.RegisterListFunc(".zip", List)
}
