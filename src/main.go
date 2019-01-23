package main

import (
	"fmt"
)

func main() {

	paths := readFolder("/Go_parser_core/texts/")
	files := initFileObjects(paths)

	for _, file := range files {
		fmt.Printf("File key: %s, file hash: %s, file path: %s, file size: %d\n",
			file.fileUniqueKey, file.fileHash, file.filePath, file.fileSize)
	}

	idx := indexing{&files}

	for _, file := range *idx.filesToIndex {
		fmt.Printf("I'm file to index: %s\n", file)
	}

}