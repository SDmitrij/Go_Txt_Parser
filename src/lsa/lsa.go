package lsa

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/big"
)

const ( dimToRender int = 3
		N           int = 2 )

type LatentSemanticAnalysis struct {
	Indexer Indexing
	Fm      *FrequencyMatrix
}

type FrequencyMatrix struct {
	frequencyMatrixVectors	*[][]int
	termsPerFile 			 []int
	tFIdf 				    *[][]float64
	SVD                     *singularValueDecomposition
	uniqueTerms				*[]string
	cosSimilarityDocs       *[][]float64
}

type singularValueDecomposition struct {
	U				mat.Matrix
	V				mat.Matrix
	S				[]float64
	fm				*FrequencyMatrix
	dataToRender	*map[string][]float64
}

func (lsa *LatentSemanticAnalysis) InvokeLsa() {

	fm                      := lsa.setFrequencyMatrix()
	fm.tFIdf                 = fm.setTfIdf()
	fm.SVD                   = fm.setSingularValueDecomposition(true)
	lsa.Fm                   = fm
	fm.SVD.dataToRender      = fm.SVD.prepareSvdDataToRender()
	fm.cosSimilarityDocs 	 = fm.SVD.cosineSimilarityOfDocuments()

	// Plotting
	fm.SVD.createHistSvdSPlot()
	fm.SVD.createTermDocumentDependencyPlot(&lsa.Indexer.Files)
	fm.SVD.cosineSimDocsPlot(&lsa.Indexer.Files)
}

/**
Set frequency matrix
TODO remove slices to consume memory usage
 */
func (lsa *LatentSemanticAnalysis) setFrequencyMatrix() *FrequencyMatrix {

	var fMatrix      				[][]int
	var termsPerFile 				[]int
	filesTerms, uniqueFilesTerms := lsa.Indexer.getTheWholeListOfTerms()

	// Fill array of vectors that contains term's matches in docs
	createTermFrequencyVec := func (matchTerm string) []int {
		vector := make([]int, len(*filesTerms))

		for i, fileTerm := range *filesTerms {
			termsPerFile = append(termsPerFile, len(fileTerm))
			for _, term := range fileTerm {
				if matchTerm == term {
					vector[i]++
				}
			}
		}

		return vector
	}

	for _, uniqueTerm := range *uniqueFilesTerms {
		vector := createTermFrequencyVec(uniqueTerm)
		var matcher int
		for _, elem := range vector {
			if elem != 0 {
				matcher++
			}
		}
		if matcher != 1 {
			fMatrix = append(fMatrix, vector)
		}
	}

	return &FrequencyMatrix{&fMatrix, termsPerFile, new([][]float64),
		       &singularValueDecomposition{}, uniqueFilesTerms, new([][]float64)}
}

/**
Set Term Frequency – Inverse Document Frequency matrix
TODO remove slices to consume memory usage
 */
func (fm *FrequencyMatrix) setTfIdf() *[][]float64 {

	var tFIdfMat [][]float64

	vecToTfIdf := func (vector []int) []float64 {
		tfIdfVector := make([]float64 , len(vector))
		nonZeroColumns := 0

		for _, elem := range vector {
			if elem != 0 {
				nonZeroColumns++
			}
		}
		for i := 0; i < len(vector); i++ {
			if vector[i] != 0 {
				tfIdfVector[i] =
					(float64(vector[i]) / float64(fm.termsPerFile[i])) * math.Log(float64(len(fm.termsPerFile)) / float64(nonZeroColumns))
			}
		}

		return tfIdfVector
	}

	for _, vector := range *fm.frequencyMatrixVectors {
		tFIdfMat = append(tFIdfMat, vecToTfIdf(vector))
	}

	return &tFIdfMat
}

/**
Set singular value decomposition according to term frequency – inverse document frequency matrix
 */
func (fm *FrequencyMatrix) setSingularValueDecomposition(print bool) *singularValueDecomposition {

	// Matrix print
	matPrint := func (X mat.Matrix) {
		fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
		fmt.Printf("%v\n", fa)
	}

	var nDim, mDim, it	int
	var minNM			float64
	var S				[]float64
	nDim                = len(*fm.tFIdf)
	mDim                = len((*fm.tFIdf)[0])
	toSVDVec           := make([]float64, nDim * mDim)

	for _, vector := range *fm.tFIdf {
		for _, elem := range vector {
			toSVDVec[it] = elem
			it++
		}
	}

	SVD         := mat.SVD{}
	toSVD       := mat.NewDense(nDim, mDim, toSVDVec)
	minNM        = math.Min(float64(nDim), float64(mDim))

	SVD.Factorize(toSVD, mat.SVDThin)

	// Left singular vector and right singular vector
	U := mat.NewDense(nDim, int(minNM), make([]float64, nDim * int(minNM)))
	V := mat.NewDense(mDim, int(minNM), make([]float64, mDim * int(minNM)))

	// Extract left singular vector and right singular vector, singular values
	S = SVD.Values(S)
	SVD.VTo(V)
	SVD.UTo(U)

	if print {
		// Print results
		fmt.Println(mDim, nDim)
		fmt.Println("U:")
		matPrint(U)
		fmt.Println("V:")
		matPrint(V)
		fmt.Println("S:")
		fmt.Println(S)
	}

	return &singularValueDecomposition{U, V, S, fm, new(map[string][]float64)}
}

/**
Extract data to render lsa plot
 */
func (svd *singularValueDecomposition) prepareSvdDataToRender() *map[string][]float64 {

	dataToRender := make(map[string][]float64)
	dimsToRender := make([]int, dimToRender)

	// Need to extract the most important two dimensions to draw the lsa plots
	setDimImportanceBySingularValues := func() []float64 {
		// Frequency analysis params
		k := len(svd.S)

		getMinMaxElem := func() (float64, float64) {
			min, max := svd.S[0], svd.S[0]
			for _, e := range svd.S {
				if e < min {
					min = e
				}
				if e > max {
					max = e
				}
			}

			return min, max
		}

		minS, maxS           := getMinMaxElem()
		spread               := maxS - minS
		h                    := spread / float64(k)
		intervals            := make([][]float64, k)
		relativeFrequency    := make([]float64, k)

		for i := range intervals {
			intervals[i] = make([]float64, N)
		}

		current := minS
		for i := 0; i < k; i++ {
			for j := 0; j < N; j++ {
				if j == 0 {
					intervals[i][j] = current
				} else {
					intervals[i][j] = current + h
					current = intervals[i][j]
				}
			}
		}

		for _, s := range svd.S {
			for i := 0; i < k; i++ {
				// Compare according to float numbers
				sValToCmp   := big.NewFloat(s)
				firstToCmp  := big.NewFloat(intervals[i][0])
				secondToCmp := big.NewFloat(intervals[i][1])

				if sValToCmp.Cmp(firstToCmp) >= 0 && sValToCmp.Cmp(secondToCmp) <= 0 {
					relativeFrequency[i] += 1
					relativeFrequency[i] = relativeFrequency[i] / float64(k)
				}
			}
		}

		return relativeFrequency
	}

	svdSValImportance := setDimImportanceBySingularValues()
	tmp := make([]float64, len(svdSValImportance))
	copy(tmp, svdSValImportance)

	for i := 0; i < dimToRender; i++ {
		for key, s := range tmp {
			if big.NewFloat(s).Cmp(big.NewFloat(tmp[0])) > 0 {
				dimsToRender[i] = key
			}
		}
		tmp[dimsToRender[i]] = float64(0)
	}

	// We throw out first dimension cause' we do not center the matrix
	firstDim, secondDim := dimsToRender[1], dimsToRender[2]
	r, c := svd.U.Dims()
	firstDimColU, secondDimColU := make([]float64, r), make([]float64, r)
	mat.Col(firstDimColU, firstDim, svd.U)
	mat.Col(secondDimColU, secondDim, svd.U)

	firstDimRowV, secondDimRowV := make([]float64, c), make([]float64, c)
	mat.Row(firstDimRowV, firstDim, svd.V)
	mat.Row(secondDimRowV, secondDim, svd.V)

	// Data to plot
	dataToRender["S_TO_HIST"] = svdSValImportance
	// Points
	dataToRender["U_TO_X"] = firstDimColU
	dataToRender["U_TO_Y"] = secondDimColU
	dataToRender["V_TO_X"] = firstDimRowV
	dataToRender["V_TO_Y"] = secondDimRowV

	return &dataToRender
}

/**
Calculate cos similarity between documents
 */
func (svd *singularValueDecomposition) cosineSimilarityOfDocuments() *[][]float64 {

	vectorsToSim := make([]mat.VecDense, len((*svd.dataToRender)["V_TO_X"]))
	cosSimValues := make([][]float64, 	 len((*svd.dataToRender)["V_TO_X"]))

	for i := 0; i < len((*svd.dataToRender)["V_TO_X"]); i++ {
		vectorsToSim[i] = *mat.NewVecDense(2,
			[]float64{ (*svd.dataToRender)["V_TO_X"][i], (*svd.dataToRender)["V_TO_Y"][i] })
	}

	for toCmp := 0; toCmp < len(vectorsToSim); toCmp++ {

		tmpValues := make([]float64, len(vectorsToSim))

		for vec := 0; vec < len(vectorsToSim); vec++ {

			if vec != toCmp {
				//  cos(Theta) = (A x B) / (sqrt(sum(A^2))|norm(A) x sqrt(sum(B^2)|norm(B)))
				tmpValues[vec] = (mat.Dot(&vectorsToSim[toCmp], &vectorsToSim[vec])) /
					( /*norm(A)*/ math.Sqrt(mat.Dot(&vectorsToSim[toCmp], &vectorsToSim[toCmp])) *
					  /*norm(B)*/ math.Sqrt(mat.Dot(&vectorsToSim[vec], &vectorsToSim[vec])) )
			}
		}

		cosSimValues[toCmp] = tmpValues
	}

	return &cosSimValues
}



