package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/big"
)

type latentSemanticAnalysis struct {
	files []File
	indexer indexing
	fm *frequencyMatrix
}

type frequencyMatrix struct {
	frequencyMatrixVectors *[][]int
	termsPerFile []int
	tFIdf *[][]float64
	SVD singularValueDecomposition
	uniqueTerms *[]string
	lsa *latentSemanticAnalysis
}

type singularValueDecomposition struct {
	U mat.Matrix
	V mat.Matrix
	S []float64
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := lsa.setFrequencyMatrix()
	fm.tFIdf = fm.setTfIdf()
	fm.SVD = *fm.setSingularValueDecomposition(true)
	lsa.fm = fm
	fm.SVD.prepareSvdDataToRender()
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix() *frequencyMatrix {

	var fMatrix [][]int
	var termsPerFile []int
	filesTerms, uniqueFilesTerms := lsa.indexer.getTheWholeListOfTerms()

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

	return &frequencyMatrix{&fMatrix, termsPerFile, &[][]float64{},
		singularValueDecomposition{}, uniqueFilesTerms, lsa }
}

/**
Set Term Frequency – Inverse Document Frequency matrix
 */
func (fm *frequencyMatrix) setTfIdf() *[][]float64 {

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
func (fm *frequencyMatrix) setSingularValueDecomposition(print bool) *singularValueDecomposition {

	// Matrix print
	matPrint := func (X mat.Matrix) {
		fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
		fmt.Printf("%v\n", fa)
	}

	var nDim, mDim, it int
	var minNM float64
	var S []float64
	nDim = len(*fm.tFIdf)
	mDim = len((*fm.tFIdf)[0])
	toSVDVec := make([]float64, nDim * mDim)

	for _, vector := range *fm.tFIdf {
		for _, elem := range vector {
			toSVDVec[it] = elem
			it++
		}
	}

	SVD := mat.SVD{}
	toSVD := mat.NewDense(nDim, mDim, toSVDVec)

	minNM = math.Min(float64(nDim), float64(mDim))
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

	return &singularValueDecomposition{U, V, S}
}

/**
Extract data to render lsa plot
 */
func (svd *singularValueDecomposition) prepareSvdDataToRender() {

	// Need to extract the most important two dimensions to draw the lsa plots
	setDimImportanceBySingularValues := func() []float64 {
		// Frequency analysis params
		const n = 2
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

		minS, maxS := getMinMaxElem()
		spread := maxS - minS
		h := spread / float64(k)
		intervals := make([][]float64, k)
		relativeFrequency := make([]float64, k)

		for i := range intervals {
			intervals[i] = make([]float64, n)
		}

		current := minS
		for i := 0; i < k; i++ {
			for j := 0; j < n; j++ {
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
				sValToCmp := big.NewFloat(s)
				firstToCmp := big.NewFloat(intervals[i][0])
				secondToCmp := big.NewFloat(intervals[i][1])

				if sValToCmp.Cmp(firstToCmp) >= 0 && sValToCmp.Cmp(secondToCmp) <= 0 {
					relativeFrequency[i] += 1
					relativeFrequency[i] = relativeFrequency[i] / float64(k)
				}
			}
		}

		return relativeFrequency
	}
}



