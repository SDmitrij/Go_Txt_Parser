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
func (idx *indexing) indexing() {
	for _, files := range *idx.filesToIndex {
		fmt.Println(files)
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