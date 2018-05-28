package kvs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	items    map[string]*item
	dir      string
	mutex    sync.Mutex
	MaxItems int
	MaxSize  int64
}

func (f *File) Default() error {
	dir, err := ioutil.TempDir("", "dewy")
	if err != nil {
		return err
	}
	f.dir = dir
	f.MaxSize = 64 * 1024 * 1024
	return nil
}

func (f *File) Read(key string) (*item, error) {
	return &item{}, nil
}

func (f *File) Write(key string, data []byte) error {
	dirstat, err := os.Stat(f.dir)
	if err != nil {
		return err
	}

	if !dirstat.Mode().IsDir() {
		return errors.New("File.dir is not dir")
	}
	if dirstat.Size() > f.MaxSize {
		return errors.New("Max size has been reached")
	}

	p := filepath.Join(f.dir, key)
	if isFileExist(p) {
		return errors.New(fmt.Sprintf("file already exists: %s", p))
	}

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)

	return nil
}

func (f *File) Delete(key string) bool {
	return true
}

func (f *File) List() []string {
	return []string{""}
}

func isFileExist(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}
