package main

import (
	"fmt"
	"regexp"
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
	re := regexp.MustCompile(`[A-Za-z']+|[!?,.]`)
	for _, file := range *idx.filesToIndex {
		fileStrings, err := file.getAllStringsOfFile(file.filePath)
		if err != nil {
			fmt.Printf("There is an error in strings reader on file: %s\n", file.filePath)
		}
		idx.filesRepo.createTableStrings(file.fileUniqueKey)
		idx.filesRepo.createTableWordsElem(file.fileUniqueKey)
		fmt.Printf("Strings of file %s\n", file.filePath)
		for lineCounter, strFile := range *fileStrings {
			fmt.Printf("%d. %s", lineCounter, strFile)
			var toStringRepo = map[string] string {"file_key": file.fileUniqueKey, "str_of_file": strFile}
			idx.filesRepo.insIntoTableStrings(toStringRepo, lineCounter)
			wordsElem := re.FindAllString(strFile, -1)
			for _, wordElem := range wordsElem {
				var toWrdElemRepo = map[string] string {"file_key": file.fileUniqueKey, "wrd_elem_of_file": wordElem}
				idx.filesRepo.insIntoTableWordsElem(toWrdElemRepo, lineCounter)
			}
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