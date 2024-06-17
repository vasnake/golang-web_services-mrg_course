package blobstorage

import (
	"io"
	"os"
)

type FSStorage struct {
	path string
}

func NewFSStorage(path string) (*FSStorage, error) {
	return &FSStorage{
		path: path,
	}, nil
}

func (st *FSStorage) Put(data io.ReadSeeker, objectName, contentType string, userID uint32) error {
	newFile, err := os.Create(st.path + objectName)
	if err != nil {
		return err
	}
	_, err = io.Copy(newFile, data)
	if err != nil {
		return err
	}
	newFile.Sync()
	newFile.Close()
	return nil
}
