package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	filePath string
	fileUniqueKey string
	fileHash string
	fileSize int
}

func getMd5HashOfFile(filePath string) (string, error){
	var retMd5Value string

	file, err := os.Open(filePath)
	if err != nil {
		return retMd5Value, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return retMd5Value, err
	}
	hashInBytes := hash.Sum(nil)[:32]
	retMd5Value = hex.EncodeToString(hashInBytes)
	return retMd5Value, err
}

func getMd5FileUniqueKey(filename string) (string, error){


}

func initFileObjects(filePaths []string) []File{
	var files []File
	for _, path := range filePaths {
		filesHash, err := getMd5HashOfFile(path)
		if err != nil {
			files = append(files, File{filePath:path, fileHash:filesHash})
		}

	}
}

func readFolder(dir string) []string {
	var filesPaths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		filesPaths = append(filesPaths, path)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return filesPaths
}
