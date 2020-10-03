package archive

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
)

func List(name string) ([]FileInfo, error) {
	ext := path.Ext(name)
	fn := findListFn(ext)
	if fn == nil {
		return []FileInfo{}, fmt.Errorf("unsupported file type: %s", ext)
	}
	return fn(name)
}

func findListFn(ext string) ListFunc {
	fns := atomicListFns.Load().([]extFunc)
	for _, f := range fns {
		if ext == f.ext {
			return f.fn
		}
	}
	return nil
}

type FileInfo struct {
	os.FileInfo
	Path string
}

type ListFunc func(string) ([]FileInfo, error)

type extFunc struct {
	ext string
	fn  ListFunc
}

var (
	listMu        sync.Mutex
	atomicListFns atomic.Value
)

func RegisterListFunc(ext string, list ListFunc) {
	listMu.Lock()
	fns := atomicListFns.Load().([]extFunc)
	atomicListFns.Store(append(fns, extFunc{ext, list}))
	listMu.Unlock()
}

func init() {
	listMu.Lock()
	atomicListFns.Store([]extFunc{})
	listMu.Unlock()
}
