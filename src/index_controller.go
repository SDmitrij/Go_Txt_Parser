package main

import (
	"fmt"
	"regexp"
	"strings"
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
func (idx *indexing) trueIndexing() {
	re := regexp.MustCompile(`[A-Za-z']+|\b`)
	for _, file := range *idx.filesToIndex {
		// Get all strings of current file
		fileStrings := file.getAllStringsOfFile(file.filePath)
		idx.filesRepo.createTableStrings(file.fileUniqueKey, "tbl_str_pref")
		idx.filesRepo.createTableWords(file.fileUniqueKey, "tbl_wrd_pref")
		fmt.Printf("Strings of file %s\n", file.filePath)
		// Index strings of file
		for lineCounter, strFile := range *fileStrings {
			fmt.Printf("%d. %s", lineCounter, strFile)
			var toStringRepo = map[string] string {"file_key": file.fileUniqueKey, "str_of_file": strFile}
			idx.filesRepo.insIntoTableStrings(toStringRepo, "tbl_str_pref", lineCounter)
			wordsElem := re.FindAllString(strFile, -1)
			// Index words and elements of string
			for _, wordElem := range wordsElem {
				var toWrdRepo = map[string] string {"file_key": file.fileUniqueKey, "wrd_of_file": wordElem}
				idx.filesRepo.insIntoTableWords(toWrdRepo, "tbl_wrd_pref", lineCounter)
			}
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

//TODO words to compare must be cased
func (idx *indexing) removeStopSymbols(stringOfFile []string) []string {
	// Simple list of english stop words
	stopWords := []string {"i", "i'm", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your", "yours",
		"yourself", "yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its", "itself",
		"they", "them", "their", "theirs", "themselves", "what", "which", "who", "whom", "this", "that", "these",
		"those", "am", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "having", "do", "does",
		"did", "doing", "a", "an", "the", "and", "but", "if", "or", "because", "as", "until", "while", "of", "at", "by",
		"for", "with", "about", "against", "between", "into", "through", "during", "before", "after", "above", "below",
		"to", "from", "up", "down", "in", "out", "on", "off", "over", "under", "again", "further", "then", "once",
		"here", "there", "when", "where", "why", "how", "all", "any", "both", "each", "few", "more", "most", "other",
		"some", "such", "no", "nor", "not", "only", "own", "same", "so", "than", "too", "very", "can", "will",
		"just", "don't", "should", "now"}

	// Anon. func that return difference between two arrays
	var difference func(a, b []string) []string
	// Anon. func to convert words in array in lower case
	var strArrToLow func(a []string) []string

	strArrToLow = func(a []string) []string {
		var res []string
		for _, value := range a {
			res = append(res, strings.ToLower(value))
		}

		return res
	}

	difference = func(a, b []string) []string {
		mb := map[string]bool{}
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

	lowStrArr := strArrToLow(stringOfFile)
	strDiffer := difference(lowStrArr, stopWords)

	return strDiffer
}