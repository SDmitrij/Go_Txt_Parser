package main

import (
	"fmt"
)

func main() {

	paths := getMainFilesInfo("/Go_parser_core/texts/")
	files := initFileObjects(paths)

	for _, file := range files {
		fmt.Printf("File key: %s, file hash: %s, file path: %s, file size: %d\n",
			file.fileUniqueKey, file.fileHash, file.filePath, file.fileSize)
	}

	idx := indexing{&files}
	fmt.Println(*idx.filesToIndex)

}