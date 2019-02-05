package main

import (
	"fmt"
)

/**
This struct describes indexing process
 */
type indexing struct {
	filesToIndex *[]File
	filesRepo *filesRepo
}

/**
Init files main info into table
 */
func (idx *indexing) initFilesInfo() {
	for _, file := range *idx.filesToIndex {
		idx.filesRepo.insIntoMainInfoFileTable(file)
	}
}

/**
Indexing current directory files
 */
func (idx *indexing) indexing(){
	for _, file := range *idx.filesToIndex {
		fileStrings, err := file.getAllStringsOfFile(file.filePath)
		if err != nil {
			fmt.Printf("There is an error in strings reader on file: %s\n", file.filePath)
		}

		fmt.Printf("Strings of file %s\n", file.filePath)
		for i, strFile := range *fileStrings {
			fmt.Printf("%d. %s", i, strFile)

		}

	}
}

/**
Search all matches and get results
 */
func (idx *indexing) searching() {

}

/**
Exclude or include files to index
 */
func (idx *indexing) excludeOrIncludeFiles() {

}