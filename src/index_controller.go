package main

import (
	"fmt"
	"github.com/caneroj1/stemmer"
	"strings"
	"unicode"
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
		// Get previous file data
		prevFileData := idx.filesRepo.getFileInfoAsObj(file.fileUniqueKey)
		if (File{}) != prevFileData {
			if prevFileData.fileHash != file.fileHash && prevFileData.fileSize != file.fileSize {
				idx.filesRepo.insIntoMainInfoFileTable(file)
				idx.trueIndexing(file)
			}
		} else {
			idx.filesRepo.insIntoMainInfoFileTable(file)
			idx.trueIndexing(file)
		}

		fmt.Printf("File key: %s, file hash: %s, file path: %s, file size: %d\n",
			file.fileUniqueKey, file.fileHash, file.filePath, file.fileSize)
	}
}

/**
Indexing current directory files
 */
func (idx *indexing) trueIndexing(file File) {
	// Get all strings of current file
	fileStrings := file.getAllStringsOfFile(file.filePath)
	idx.filesRepo.createTableStrings(file.fileUniqueKey, "tbl_str_pref")
	idx.filesRepo.createTableWords(file.fileUniqueKey, "tbl_wrd_pref")

	// Index strings of file
	for lineCounter, strFile := range *fileStrings {
		var toStringRepo = map[string] string {"file_key": file.fileUniqueKey, "str_of_file": strFile}
		idx.filesRepo.insIntoTableStrings(toStringRepo, "tbl_str_pref", lineCounter)
		wordsElemStop := idx.removeStopSymbols(strFile)
		for _, wordElem := range wordsElemStop {
			stemStopWrd := stemmer.Stem(wordElem)
			var toWrdRepo = map[string] string {"file_key": file.fileUniqueKey, "wrd_of_file": strings.ToLower(stemStopWrd)}
			idx.filesRepo.insIntoTableWords(toWrdRepo, "tbl_wrd_pref", lineCounter)
		}
	}
}

func (idx *indexing) lsaIndexing() {

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

func (idx *indexing) removeStopSymbols(stringOfFile string) []string {
	// Simple list of english stop words
	stopWords := []string {"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your",
		"yours", "yourself", "yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its",
		"itself", "they", "them", "their", "theirs", "themselves", "what", "which", "who", "whom", "this", "that", "these",
		"those", "am", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "having", "do", "does",
		"did", "doing", "a", "an", "the", "and", "but", "if", "or", "because", "as", "until", "while", "of", "at", "by",
		"for", "with", "about", "against", "between", "into", "through", "during", "before", "after", "above", "below",
		"to", "from", "up", "down", "in", "out", "on", "off", "over", "under", "again", "further", "then", "once",
		"here", "there", "when", "where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other",
		"some", "such", "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "can", "will",
		"just", "don't", "should", "now", "m", "ll", "d", "s", "t"}

	// Anon. func that return difference between two arrays
	var difference func(a, b []string) []string

	// Fast array diff
	difference = func(a, b []string) []string {
		mb := make(map[string]bool)
		var ab []string
		for _, x := range b {
			mb[x] = true
		}
		for _, x := range a {
			if _, ok := mb[x]; !ok {
				ab = append(ab, x)
			}
		}

		return ab
	}

	toLow := strings.ToLower(stringOfFile)
	words := strings.FieldsFunc(toLow, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	differ := difference(words, stopWords)

	return differ
}