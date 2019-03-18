package main

type latentSemanticAnalysis struct {
	files *[] File
	indexer *indexing
}

type frequencyMatrix struct {
	frequencyMatrixVectors map[string][]int
	tFIdf map[string][]int
	lsa latentSemanticAnalysis
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := frequencyMatrix{*lsa.setFrequencyMatrix(true), make(map[string][]int),
		*lsa }
	fm.setTfIdf()
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix(lessMatch bool) *map[string][]int {

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

	return &termVectors
}

func (fm *frequencyMatrix) setTfIdf() {

	vecToTfIdf := func(vector *[]int) {

	}

	for _, fileTerms := range fm.frequencyMatrixVectors {
		for _, termsVector := range fileTerms {

		}
	}
}



