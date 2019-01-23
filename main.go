package main

import (
	"fmt"
)

func main() {

	paths := readFolder("/GoProjects/texts/")
	files := initFileObjects(paths)

	for _, file := range files {
		fmt.Printf("File key: %s, file hash: %s\n", file.fileUniqueKey, file.fileHash)
	}

}