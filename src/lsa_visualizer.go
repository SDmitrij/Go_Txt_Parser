package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (svd *singularValueDecomposition) createHistSvdSPlot(data []float64) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Singular values histogram"
	p.X.Label.Text = "Singular values"
	p.Y.Label.Text = "Importance"

	v := make(plotter.Values, len(data))

	for i := range v {
		v[i] = data[i]
	}

	h, err := plotter.NewHist(v, len(data))
	if err != nil {
		panic(err)
	}

	h.Normalize(1)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "svd_s_val_importance.png"); err != nil {
		panic(err)
	}
}

func (svd *singularValueDecomposition) createTermDocumentDependencyPlot() {

}
