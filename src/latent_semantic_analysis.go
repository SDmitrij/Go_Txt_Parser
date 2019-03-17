package main

type latentSemanticAnalysis struct {
	files *[] File
	repository *filesRepo
	indexer *indexing
}

func (lsa *latentSemanticAnalysis) setFrequencyMatrix() *map[string][]int{
	filesTerms := lsa.indexer.getTheWholeListOfTerms()
	termVectors := make(map[string][]int)

	searchMatches := func (matchTerm string) {
		tmpMatchVec := make([]int, len(*filesTerms))

		for it, fileTerm := range *filesTerms {
			for _, term := range fileTerm {
				if matchTerm == term {
					tmpMatchVec[it] += 1
				}
			}
		}

		termVectors[matchTerm] = tmpMatchVec
	}

	for _, fileTerm := range *filesTerms {
		for _, term := range fileTerm {
			searchMatches(term)
		}
	}

	return &termVectors
}



