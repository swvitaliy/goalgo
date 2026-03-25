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
	// Данные из бенчмарков (время в секундах на 10M вызовов)
	hrwNodes := []float64{3, 5, 10, 20, 50, 100}
	hrwTimes := []float64{0.0007897, 1.031237428, 1.130151591, 1.015335613, 1.725062567, 2.3182724}

	wrhNodes := []float64{3, 5, 10, 20, 50, 100}
	wrhTimes := []float64{1.560922831, 2.131583121, 3.530437827, 5.979558726, 12.874355719, 24.087611798}

	consistentNodes := []float64{3, 5, 10, 20, 50, 100}
	consistentTimes := []float64{0.083641596, 1.06038645, 1.041598835, 1.096367033, 1.221272036, 1.324827715}

	// Построение графика для HRW
	p1 := plot.New()
	p1.Title.Text = "HRW Execution Time vs Number of Nodes"
	p1.X.Label.Text = "Number of Nodes"
	p1.Y.Label.Text = "Time (seconds)"

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
	p2.Y.Label.Text = "Time (seconds)"

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
	p3.Y.Label.Text = "Time (seconds)"

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
