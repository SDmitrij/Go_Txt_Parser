package lsa

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"math"
	"math/rand"
	"path/filepath"
	"time"
)

func createPoints(x, y []float64) plotter.XYs {

	points := make(plotter.XYs, len(x))
	for i := range points {
		points[i].X = x[i]
		points[i].Y = y[i]
	}

	return points
}

func createLabelDocsPlot(docPoints plotter.XYs, files *[]File) (plotter.Scatter, plotter.Labels) {

	d, err := plotter.NewScatter(docPoints)
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

	l, err := plotter.NewLabels(plotter.XYLabels{docPoints, FileNames })
	if err != nil {
		panic(err)
	}
	l.XOffset = 2 * vg.Millimeter

	return *d, *l
}

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
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "svd_s_val_importance.png"); err != nil {
		panic(err)
	}
}

func (svd *singularValueDecomposition) createTermDocumentDependencyPlot(files *[]File) {

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

	d, l := createLabelDocsPlot(VPoints, files)

	// Add objects to plot
	p.Add(t, &d, &l)

	p.Legend.Padding = 2 * vg.Millimeter
	p.Legend.Add("term", t)
	p.Legend.Add("document", &d)

	// Save the plot to a PNG file.
	if err := p.Save(8 * vg.Inch, 8 * vg.Inch, "term_doc_dependence.png"); err != nil {
		panic(err)
	}
}

func (svd *singularValueDecomposition) cosineSimDocsPlot(files *[]File) {

	var randUint   func(min, max int) uint8
		rand.Seed(time.Now().Unix())
		randUint = func(min, max int) uint8 {
			return uint8(rand.Intn(max - min) + min)
		}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text        = "Cosine similarity"
	p.X.Label.Text      = "First dimension"
	p.Y.Label.Text      = "Second dimension"
	VPoints 		    := createPoints((*svd.dataToRender)["V_TO_X"], (*svd.dataToRender)["V_TO_Y"])

	d, label := createLabelDocsPlot(VPoints, files)
	p.Add(&d, &label)

	cosineSimilarityValues := *svd.fm.cosSimilarityDocs

	for i := 0; i < len(cosineSimilarityValues); i++ {

		var minIdxOnSub int

		for j := 0; j < len(cosineSimilarityValues[i]); j++ {
			if math.Abs(cosineSimilarityValues[i][j]) < math.Abs(cosineSimilarityValues[i][0]) &&
				cosineSimilarityValues[i][j] != 0 {
				minIdxOnSub = j
			}
		}

		lineColor := color.RGBA{
			R: randUint(0, 255),
			G: randUint(0, 255),
			B: randUint(0, 255),
			A: 255 }

		line, err := plotter.NewLine(plotter.XYs{ VPoints[i], VPoints[minIdxOnSub] })
		if err != nil {
			panic(err)
		}
		line.Color = lineColor
		p.Add(line)
		p.Legend.Add(fmt.Sprintf("sim = %f", math.Abs(cosineSimilarityValues[i][minIdxOnSub])), line)
	}

	// Save the plot to a PNG file.
	if err := p.Save(8 * vg.Inch, 8 * vg.Inch, "cosine_similarity_docs.png"); err != nil {
		panic(err)
	}
}
