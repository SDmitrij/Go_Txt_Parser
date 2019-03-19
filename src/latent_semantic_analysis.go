package main

import "math"

type latentSemanticAnalysis struct {
	files *[] File
	indexer *indexing
	fm *frequencyMatrix
}

type frequencyMatrix struct {
	frequencyMatrixVectors *map[string][]int
	wordsPerDoc *map[string]int
	tFIdf map[string][]int
	lsa latentSemanticAnalysis
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := lsa.setFrequencyMatrix(true)
	lsa.fm = &fm
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix(lessMatch bool) frequencyMatrix {

	wordsPerDoc := make(map[string]int)
	filesTerms := lsa.indexer.getTheWholeListOfTerms()
	termVectors := make(map[string][]int)

	// Fill array of vectors that contains term's matches in docs
	searchMatchesVectors := func (matchTerm string) {
		tmpMatchVec := make([]int, len(*filesTerms))
		var iter int
		for _, fileTerm := range *filesTerms {
			for _, term := range fileTerm {
				if matchTerm == term {
					tmpMatchVec[iter]++
				}
			}
			iter++
		}

		termVectors[matchTerm] = tmpMatchVec
	}

	// Go through array of terms
	for filename, fileTerm := range *filesTerms {
		wordsPerDoc[filename] = len(fileTerm)
		for _, term := range fileTerm {
			searchMatchesVectors(term)
		}
	}

	// Remove terms that match less that two times
	if lessMatch {
		var unsetVectors []string
		for vector, termVector := range termVectors {
			var matcher int
			for _, termMatch := range termVector {
				if termMatch != 0 {
					matcher++
				}
			}
			if matcher == 1 {
				unsetVectors = append(unsetVectors, vector)
			}
		}
		for _, vecKey := range unsetVectors {
			delete(termVectors, vecKey)
		}
	}

	return frequencyMatrix{&termVectors, &wordsPerDoc, map[string][]int{},
		*lsa}
}

func (fm *frequencyMatrix) setTfIdf() {

	vecToTfIdf := func(arr *[]int, filename string) {

		vector := *arr
		tfIdfVector := make([]int , len(vector))
		var nonZeroColumns int

		for _, elem := range vector {
			if elem != 0 {
				nonZeroColumns++
			}
		}

		for i := 0; i < len(vector); i++ {
			tfIdfVector[i] =
				(vector[i] / fm.wordsPerDoc[filename]) * math.Log(float64(len(fm.wordsPerDoc) / nonZeroColumns))
		}
	}

	for filename, fileTerms := range *fm.frequencyMatrixVectors {
		for _, termsVector := range fileTerms {

		}
	}
}



