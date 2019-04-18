package lsa

import (
	"bufio"
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
	filePath      string
	fileUniqueKey string
	fileHash 	  string
	fileSize      int
}

/**
Get md 5 hash of each file
 */
func getMdFiveHashOfFile(filePath string) (string, error){
	var mdFiveHash string
	file, err := os.Open(filePath)
	if err != nil {
		return mdFiveHash, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return mdFiveHash, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	mdFiveHash = hex.EncodeToString(hashInBytes)

	return mdFiveHash, err
}

/**
Get md 5 unique key of each file
 */
func getMdFiveFileUniqueKey(filename string) string {
	var mdFiveUniqueKey string
	hash := md5.New()
	if _, err := io.WriteString(hash, filename); err != nil {
		return mdFiveUniqueKey
	}
	return hex.EncodeToString(hash.Sum(nil)[:16])
}

/**
Init file objects
 */
func InitFileObjects(filesInfo map[string]int64) []File {
	var files []File
	for path, size := range filesInfo {
		filesHash, errFileHash := getMdFiveHashOfFile(path)
		uniqueFileKey := getMdFiveFileUniqueKey(filepath.Base(path))
		fileSize := size
		if errFileHash == nil{
			files = append(files,
				File{filePath: path, fileUniqueKey: uniqueFileKey, fileHash: filesHash, fileSize: int(fileSize)})
		}
	}

	return files
}

/**
Read folder with texts to index and get file's size and path
 */
func GetMainFilesInfo(dir string) map[string]int64 {
	var filesInfo = make(map[string]int64)
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

/**
Get file content line by line
 */
func (f *File) getAllStringsOfFile(filePath string) *[]string {
	var lines []string
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}

	return &lines
}