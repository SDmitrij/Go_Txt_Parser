package main

import (
	"math"
)

type latentSemanticAnalysis struct {
	files *[] File
	indexer *indexing
	fm *frequencyMatrix
}

type frequencyMatrix struct {
	frequencyMatrixVectors *map[string][]int
	wordsPerDoc *[]int
	tFIdf *[][]float64
	lsa *latentSemanticAnalysis
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := lsa.setFrequencyMatrix(true)
	fm.tFIdf = fm.setTfIdf()
	lsa.fm = &fm
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix(lessMatch bool) frequencyMatrix {

	var wordsPerDoc []int
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
	for _, fileTerm := range *filesTerms {
		wordsPerDoc = append(wordsPerDoc, len(fileTerm))
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

	return frequencyMatrix{&termVectors, &wordsPerDoc, &[][]float64{},
		lsa}
}

/**
Set Term Frequency – Inverse Document Frequency matrix
 */
func (fm *frequencyMatrix) setTfIdf() *[][]float64 {

	wordsPerDoc := *fm.wordsPerDoc
	var tFIdfMatrix [][]float64

	vecToTfIdf := func(vector []int) []float64 {
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

	for _, fVector := range *fm.frequencyMatrixVectors {
		tFIdfMatrix = append(tFIdfMatrix, vecToTfIdf(fVector))
	}

	return &tFIdfMatrix
}



