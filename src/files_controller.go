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
Init file objects
 */
func initFileObjects(filesInfo map[string] int64) []File{
	var files []File
	for path, size := range filesInfo {
		filesHash, errFileHash := getMd5HashOfFile(path)
		uniqueFileKey := getMd5FileUniqueKey(filepath.Base(path))
		fileSize := size
		if errFileHash == nil{
			files = append(files, File{filePath: path, fileUniqueKey: uniqueFileKey, fileHash: filesHash, fileSize: fileSize})
		}
	}

	return files
}

/**
Read folder with texts to index and get file's size and path
 */
func getMainFilesInfo(dir string) map[string] int64 {
	var filesInfo  = make(map[string] int64)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesInfo[path] = info.Size()
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return filesInfo
}
