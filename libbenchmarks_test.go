package main

import "testing"

func TestGraphics(t *testing.T) {
	SetOptions()

	data := Data{Durations: []int64{
		6887, 8380, 9973, 11883, 14536, 15490, 16788, 17836, 18841, 19586, 21094, 21642, 22823,
		24088, 24382, 25053, 26204, 28545, 29191, 27631, 28112, 28873, 30387, 30725, 29175, 31428,
		32644, 32365, 31683, 33026, 33375, 35130, 34001, 36181, 36498, 36868, 38079, 37601, 36767,
		38309},
		Benchmark: BenchmarkModel{BenchmarkType: "ActiveRequestWithBackgroundLoad"}}
	writePlots(data)
}
