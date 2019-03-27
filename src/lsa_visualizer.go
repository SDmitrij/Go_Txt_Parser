package main

import (
	"gonum.org/v1/plot"
)

func (svd *singularValueDecomposition) createLsaPlot() {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Term - document dependency"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
}
