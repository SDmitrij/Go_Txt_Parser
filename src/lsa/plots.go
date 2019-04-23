package lsa

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"path/filepath"
)

func (svd *singularValueDecomposition) createHistSvdSPlot() {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text        = "Singular values histogram"
	p.X.Label.Text      = "Singular values"
	p.Y.Label.Text      = "Importance"

	v := make(plotter.Values, len((*svd.dataToRender)["S_TO_HIST"]))

	for i := range v {
		v[i] = (*svd.dataToRender)["S_TO_HIST"][i]
	}

	h, err := plotter.NewHist(v, len((*svd.dataToRender)["S_TO_HIST"]))
	if err != nil {
		panic(err)
	}

	h.Normalize(1)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "svd_s_val_importance.png"); err != nil {
		panic(err)
	}
}

func (svd *singularValueDecomposition) createTermDocumentDependencyPlot(files *[]File) {

	createPoints := func(x, y []float64) plotter.XYs {
		points := make(plotter.XYs, len(x))
		for i := range points {
			points[i].X = x[i]
			points[i].Y = y[i]
		}

		return points
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text        = "Term - document dependency"
	p.X.Label.Text      = "First dimension"
	p.Y.Label.Text      = "Second dimension"

	t, err := plotter.NewScatter(createPoints((*svd.dataToRender)["U_TO_X"], (*svd.dataToRender)["U_TO_Y"]))
	if err != nil {
		panic(err)
	}
	t.GlyphStyle.Shape = draw.RingGlyph{}
	t.Radius = 2 * vg.Millimeter
	t.Color = color.RGBA{R: 255, A:255}

	VPoints := createPoints((*svd.dataToRender)["V_TO_X"], (*svd.dataToRender)["V_TO_Y"])

	d, err := plotter.NewScatter(VPoints)
	if err != nil {
		panic(err)
	}
	d.GlyphStyle.Shape = draw.PyramidGlyph{}
	d.Radius = 2 * vg.Millimeter
	d.Color = color.RGBA{R:0, A:255}

	FileNames := make([]string, len(*files))
	for i, file := range *files {
		FileNames[i] = filepath.Base(file.filePath)
	}

	l, err := plotter.NewLabels(plotter.XYLabels{VPoints, FileNames })
	if err != nil {
		panic(err)
	}
	l.XOffset = 2 * vg.Millimeter

	// Add objects to plot
	p.Add(t, d, l)

	p.Legend.Padding = 2 * vg.Millimeter
	p.Legend.Add("term", t)
	p.Legend.Add("document", d)

	// Save the plot to a PNG file.
	if err := p.Save(6 * vg.Inch, 6 * vg.Inch, "term_doc_dependence.png"); err != nil {
		panic(err)
	}
}
