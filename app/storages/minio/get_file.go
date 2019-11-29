package minio

import (
	"GIG/commons"
	"github.com/minio/minio-go"
	"io"
	"os"
)

/**
Retrieve file from storage
 */
func (h Handler) GetFile(directoryName string, filename string) (*os.File, error) {
	object, err := h.Client.GetObject(directoryName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	tempDir := h.CacheDirectory + directoryName + "/"
	sourcePath := tempDir + filename

	if err = commons.EnsureDirectory(tempDir); err != nil {
		return nil, err
	}

	localFile, err := os.Create(sourcePath)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(localFile, object); err != nil {
		return nil, err
	}
	return localFile, err
}