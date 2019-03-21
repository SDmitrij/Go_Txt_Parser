package main

import (
	"math"
)

type latentSemanticAnalysis struct {
	files *[]File
	indexer *indexing
	fm *frequencyMatrix
}

type frequencyMatrix struct {
	frequencyMatrixVectors *map[string]map[string][]int
	tFIdf *map[string][]float64
	lsa *latentSemanticAnalysis
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := lsa.setFrequencyMatrix(true)
	/*
	fm.tFIdf = fm.setTfIdf()
	lsa.fm = &fm
	*/
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix(lessMatch bool) frequencyMatrix {

	// Create multidimensional map that contains file keys, terms and their frequency
	fMatrix := make(map[string]map[string][]int)
	filesTerms := lsa.indexer.getTheWholeListOfTerms()

	// Fill array of vectors that contains term's matches in docs
	searchMatchesVectors := func (matchTerm string) []int{
		matchVector := make([]int, len(*filesTerms))
		var iter int
		for _, fileTerm := range *filesTerms {
			for _, term := range fileTerm {
				if matchTerm == term {
					matchVector[iter]++
				}
			}
			iter++
		}

		return matchVector
	}

	// Go through array of terms
	for fileKey, fileTerm := range *filesTerms {
		termVectors := make(map[string][]int)
		for _, term := range fileTerm {
			termVectors[term] = searchMatchesVectors(term)
		}
		fMatrix[fileKey] = termVectors
	}

	// Remove terms that match less than two times
	if lessMatch {

		var unsetVec [][]string
		for fileUniqueKey, fileTerms := range fMatrix {
			for term, termVector := range fileTerms {
				var matcher int
				for _, vecVal := range termVector {
					if vecVal != 0 {
						matcher++
					}
				}
				if matcher == 1 {
					unsetVec = append(unsetVec, []string{fileUniqueKey, term})
				}
			}
		}

		for _, unsetInfo := range unsetVec {
			delete(fMatrix[unsetInfo[0]], unsetInfo[1])
		}
	}

	return frequencyMatrix{& fMatrix, & map[string][]float64{}, lsa }
}

/**
Set Term Frequency – Inverse Document Frequency matrix
 */
func (fm *frequencyMatrix) setTfIdf() *map[string][]float64 {

	tFIdfMatrix := make(map[string][]float64)

	vecToTfIdf := func(vector []int, fileUniqueKey string) []float64 {
		tfIdfVector := make([]float64 , len(vector))

		var nonZeroColumns int

		for _, elem := range vector {
			if elem != 0 {
				nonZeroColumns++
			}
		}

		for i := 0; i < len(vector); i++ {
			if vector[i] != 0 {
				tfIdfVector[i] =
					(float64(vector[i]) / float64(wordsPerDoc[i])) * math.Log(float64(len(wordsPerDoc)) / float64(nonZeroColumns))
			}
		}

		return tfIdfVector
	}

	for fileUniqueKey, fileTerms := range *fm.frequencyMatrixVectors {
		//tFIdfMatrix[term] = vecToTfIdf(fVector)
		for term, termVector := range fileTerms {

		}
	}

	return & tFIdfMatrix
}

func (fm *frequencyMatrix) setSingularValueDecomposition() {

}



