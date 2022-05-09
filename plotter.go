package main

import (
	"github.com/wcharczuk/go-chart"
	"log"
	"os"
)

func preparePoints(data Data) Data {
	var XValues []float64
	var YValues []float64
	for x, y := range data.Durations {
		YValues = append(YValues, float64(y))
		XValues = append(XValues, float64(x))
	}
	data.DurationX = XValues
	data.DurationY = YValues
	return data
}

func drawPlot(data Data, f *os.File) {
	plot := chart.Chart{
		Title: data.Benchmark.BenchmarkType,
		XAxis: chart.XAxis{Name: "t,sec", Style: chart.Style{Show: true}},
		YAxis: chart.YAxis{Name: "Query duration,sec", Style: chart.Style{Show: true}},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{

					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					Show: 		 true,
					Padding: chart.DefaultBackgroundPadding,
				},
				XValues: data.DurationX,
				YValues: data.DurationY,
			},
		},
	}
	err := plot.Render(chart.PNG, f)
	if err != nil {
		log.Fatal("got error while rendering plot: ", err)
	}
}

func writePlots(data Data) {
	data = preparePoints(data)
	fd, err := os.Create(myPresets.eo.FilePath + "durations_" + "plot.png")
	if err != nil {
		log.Fatal("cannot open file ", err)
	}
	drawPlot(data, fd)
	err = fd.Close()
	if err != nil {
		log.Fatal("cannot close file ", err)
	}
}
