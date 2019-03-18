package main

type latentSemanticAnalysis struct {
	files *[] File
	repository *filesRepo
	indexer *indexing
}

func (lsa *latentSemanticAnalysis) setFrequencyMatrix(lessMatch bool) *map[string][]int {

	filesTerms := lsa.indexer.getTheWholeListOfTerms()
	termVectors := make(map[string][]int)

	// Fill array of vectors that contains term's matches in docs
	searchMatchesVectors := func (matchTerm string) {
		tmpMatchVec := make([]int, len(*filesTerms))
		iter := 0
		for _, fileTerm := range *filesTerms {
			for _, term := range fileTerm {
				if matchTerm == term {
					tmpMatchVec[iter] += 1
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



