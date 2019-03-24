package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
)

type latentSemanticAnalysis struct {
	files *[]File
	indexer *indexing
	fm *frequencyMatrix
}

type frequencyMatrix struct {
	frequencyMatrixVectors *[][]int
	termsPerFile *[]int
	tFIdf *[][]float64
	lsa *latentSemanticAnalysis
}

func (lsa *latentSemanticAnalysis) invokeLsa() {
	fm := lsa.setFrequencyMatrix()
	fm.tFIdf = fm.setTfIdf()
	fm.setSingularValueDecomposition()
	lsa.fm = &fm
}

/**
Set frequency matrix
 */
func (lsa *latentSemanticAnalysis) setFrequencyMatrix() frequencyMatrix {

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

	return frequencyMatrix{&fMatrix, &termsPerFile, &[][]float64{}, lsa }
}

/**
Set Term Frequency – Inverse Document Frequency matrix
 */
func (fm *frequencyMatrix) setTfIdf() *[][]float64 {

	var tFIdfMat [][]float64

	vecToTfIdf := func(vector []int) []float64 {
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
					(float64(vector[i]) / float64((*fm.termsPerFile)[i])) * math.Log(float64(len(*fm.termsPerFile)) / float64(nonZeroColumns))
			}
		}

		return tfIdfVector
	}

	for _, vector := range *fm.frequencyMatrixVectors {
		tFIdfMat = append(tFIdfMat, vecToTfIdf(vector))
	}

	return &tFIdfMat
}

func (fm *frequencyMatrix) setSingularValueDecomposition() {

	matPrint := func(X mat.Matrix) {
		fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
		fmt.Printf("%v\n", fa)
	}

	var nDim, mDim, it int
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

	toSVD := mat.NewDense(nDim, mDim, toSVDVec)
	U := mat.NewDense(nDim, mDim, make([]float64, mDim * nDim))
	V := mat.NewDense(mDim, mDim, make([]float64, mDim * mDim))

	SVD := mat.SVD{}
	SVD.Factorize(toSVD, mat.SVDThin)
	S = SVD.Values(S)
	SVD.VTo(V)
	SVD.UTo(U)
	fmt.Println(mDim, nDim)
	fmt.Println("U:")
	matPrint(U)
	fmt.Println("V:")
	matPrint(V)
	fmt.Println("S:")
	fmt.Println(S)
}



