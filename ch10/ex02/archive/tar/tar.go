package tar

import (
	t "archive/tar"
	"io"
	"os"

	"local.gopl/archive"
)

func List(name string) (files []archive.FileInfo, err error) {
	r, err := os.Open(name)
	if err != nil {
		return
	}
	defer r.Close()

	tr := t.NewReader(r)

	for {
		h, e := tr.Next()
		if e != nil {
			if e == io.EOF {
				break
			}
			err = e
			return
		}
		files = append(files, archive.FileInfo{h.FileInfo(), h.Name})
	}
	return
}

func init() {
	archive.RegisterListFunc(".tar", List)
}
