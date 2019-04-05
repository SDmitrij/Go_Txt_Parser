package main

import (
	"github.com/caneroj1/stemmer"
	"strings"
	"unicode"
)

/**
This struct describes indexing process
 */
type indexing struct {
	filesToIndex []File
	filesRepo filesRepo
}

func (idx *indexing) invokeIndexing() {
	idx.filesRepo.initFilesRepo()
	idx.initFilesInfo()
}

/**
Init files main info into table
 */
func (idx *indexing) initFilesInfo() {

	for _, file := range idx.filesToIndex {
		// Get previous file data
		prevFileData := idx.filesRepo.getFileInfoAsObj(file.fileUniqueKey)
		if (File{}) != prevFileData {
			if prevFileData.fileHash != file.fileHash && prevFileData.fileSize != file.fileSize {
				idx.filesRepo.deleteFileInfo(file.fileUniqueKey)
				idx.filesRepo.insIntoMainInfoFileTable(file)
				idx.trueIndexing(file)
			}
		} else {
			idx.filesRepo.insIntoMainInfoFileTable(file)
			idx.trueIndexing(file)
		}
	}
}

/**
Indexing current directory files
 */
func (idx *indexing) trueIndexing(file File) {

	// Get all strings and words of current file
	fileLines := file.getAllStringsOfFile(file.filePath)
	// Get all words of current file
	fileTerms := idx.prepareWords(fileLines)
	// Create entry tables for each file
	idx.filesRepo.createTableStrings(file.fileUniqueKey, "tbl_str_pref")
	idx.filesRepo.createTableTerms(file.fileUniqueKey, "tbl_term_pref")

	// Index strings of file
	for _, strFile := range *fileLines {
		toStringRepo := map[string]string{"file_key": file.fileUniqueKey, "str_of_file": strFile}
		idx.filesRepo.insIntoTableStrings(toStringRepo, "tbl_str_pref")
	}

	// Index words
	for _, term := range *fileTerms {
		toTermRepo := map[string] string {"file_key": file.fileUniqueKey, "term_of_file": term}
		idx.filesRepo.insIntoTableTerms(toTermRepo, "tbl_term_pref")
	}
}

func (idx * indexing) removeStopSymbols(stringOfFile string) []string {

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

	// Fast array diff
    difference := func(a, b []string) []string {
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

	differ := difference(strings.FieldsFunc(strings.ToLower(stringOfFile), func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}), stopWords)

	return differ
}

func (idx *indexing) prepareWords(fileLines *[]string) *[]string {

	var prepare []string

	getStemmedWords := func(words []string) {
		iterateWords := func() {
			for _, word := range words {
				prepare = append(prepare, strings.ToLower(stemmer.Stem(word)))
			}
		}
		iterateWords()
	}

	for _, line := range *fileLines {
		getStemmedWords(idx.removeStopSymbols(line))
	}

	return &prepare
}

/**
Get the whole terms list for each file
 */
func (idx * indexing) getTheWholeListOfTerms() (*[][]string, *[]string) {

	var allFilesTerms [][]string
	var toUnique []string

	// Fast de-duplicator
	removeDuplicates := func(elements []string) []string {
		encountered := make(map[string]bool)
		var result []string
		// Create a map of all unique elements.
		for v := range elements {
			encountered[elements[v]] = true
		}
		// Place all keys from the map into a slice.
		for key := range encountered {
			result = append(result, key)
		}

		return result
	}

	for _, file := range idx.filesToIndex {
		allFilesTerms = append(allFilesTerms, *idx.filesRepo.getAllTermsOfFile(file.fileUniqueKey, "tbl_term_pref"))
	}

	for _, allFileTerm := range allFilesTerms {
		for _, term := range allFileTerm {
			toUnique = append(toUnique, term)
		}
	}

	uniqueTerms := removeDuplicates(toUnique)

	return &allFilesTerms, &uniqueTerms
}

