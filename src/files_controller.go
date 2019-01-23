package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

/**
Structure that describes file entity
 */
type File struct {
	filePath string
	fileUniqueKey string
	fileHash string
	fileSize int64
}

/**
Get md 5 hash of each file
 */
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
	hashInBytes := hash.Sum(nil)[:16]
	retMd5Value = hex.EncodeToString(hashInBytes)
	return retMd5Value, err
}

/**
Get md 5 unique key of each file
 */
func getMd5FileUniqueKey(filename string) string{
	var md5UniqueKey string
	hash := md5.New()
	if _, err := io.WriteString(hash, filename); err != nil {
		return md5UniqueKey
	}
	return hex.EncodeToString(hash.Sum(nil)[:16])
}

/**
Get file size of each file
 */
func getFileSize(filePath string) (int64, error){
	var fileSize int64
	file, err := os.Open(filePath)
	if err != nil {
		return fileSize, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return fileSize, err
	}

	fileSize = stat.Size()

	return fileSize, err
}

/**
Init file objects
 */
func initFileObjects(filePaths []string) []File{
	var files []File
	for _, path := range filePaths {
		filesHash, errFileHash := getMd5HashOfFile(path)
		uniqueFileKey := getMd5FileUniqueKey(filepath.Base(path))
		fileSize, errFileSize := getFileSize(path)
		if errFileHash == nil && errFileSize == nil {
			files = append(files, File{filePath: path, fileUniqueKey: uniqueFileKey, fileHash: filesHash, fileSize: fileSize})
		}
	}

	return files
}

/**
Read folder with texts to index
 */
func readFolder(dir string) []string {
	var filesPaths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir {
			filesPaths = append(filesPaths, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return filesPaths
}
