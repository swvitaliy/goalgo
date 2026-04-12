package main

import (
	"fmt"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Данные из бенчмарков (ns/op из bench_results_100k.txt, конвертировано в микросекунды)
	hrwNodes := []float64{3, 5, 10, 20, 50, 100}
	hrwTimes := []float64{0.7362, 0.8826, 0.7510, 1.028, 1.568, 2.443}

	wrhNodes := []float64{3, 5, 10, 20, 50, 100}
	wrhTimes := []float64{1.725, 2.068, 3.519, 5.054, 11.63, 26.24}

	consistentNodes := []float64{3, 5, 10, 20, 50, 100}
	consistentTimes := []float64{1.211, 1.091, 1.250, 1.457, 1.613, 1.667}

	// Построение графика для HRW
	p1 := plot.New()
	p1.Title.Text = "HRW Execution Time vs Number of Nodes"
	p1.X.Label.Text = "Number of Nodes"
	p1.Y.Label.Text = "Time (microseconds)"

	hrwPts := make(plotter.XYs, len(hrwNodes))
	for i := range hrwPts {
		hrwPts[i].X = hrwNodes[i]
		hrwPts[i].Y = hrwTimes[i]
	}

	err := plotutil.AddLinePoints(p1, "HRW", hrwPts)
	if err != nil {
		log.Fatal(err)
	}

	if err := p1.Save(4*vg.Inch, 4*vg.Inch, "hrw_plot.png"); err != nil {
		log.Fatal(err)
	}

	// Построение графика для WRH
	p2 := plot.New()
	p2.Title.Text = "WRH Execution Time vs Number of Nodes"
	p2.X.Label.Text = "Number of Nodes"
	p2.Y.Label.Text = "Time (microseconds)"

	wrhPts := make(plotter.XYs, len(wrhNodes))
	for i := range wrhPts {
		wrhPts[i].X = wrhNodes[i]
		wrhPts[i].Y = wrhTimes[i]
	}

	err = plotutil.AddLinePoints(p2, "WRH", wrhPts)
	if err != nil {
		log.Fatal(err)
	}

	if err := p2.Save(4*vg.Inch, 4*vg.Inch, "wrh_plot.png"); err != nil {
		log.Fatal(err)
	}

	// Построение графика для Consistent
	p3 := plot.New()
	p3.Title.Text = "Consistent Execution Time vs Number of Nodes"
	p3.X.Label.Text = "Number of Nodes"
	p3.Y.Label.Text = "Time (microseconds)"

	consistentPts := make(plotter.XYs, len(consistentNodes))
	for i := range consistentPts {
		consistentPts[i].X = consistentNodes[i]
		consistentPts[i].Y = consistentTimes[i]
	}

	err = plotutil.AddLinePoints(p3, "Consistent", consistentPts)
	if err != nil {
		log.Fatal(err)
	}

	if err := p3.Save(4*vg.Inch, 4*vg.Inch, "consistent_plot.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Graphs saved as hrw_plot.png, wrh_plot.png, consistent_plot.png")
}
