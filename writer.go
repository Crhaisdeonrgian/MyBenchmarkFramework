package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func writeIntro(file *os.File, data Data) {
	_, err := file.WriteString("BENCHMARK RESULTS:" + "\n" +
		" Benchmark type: " + data.Benchmark.BenchmarkType + "\n" +
		" Active query: " + data.Benchmark.ActiveQuery + "\n" +
		" Background query: " + data.Benchmark.BackgroundQuery + "\n")
	if err != nil {
		log.Fatal("Unable to write in file", err)
	}
}

func writeData(data Data) {
	file, err := os.Create(myPresets.eo.FilePath + myPresets.bo.driverName + ".csv")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	statfile, err := os.Create(myPresets.eo.FilePath + myPresets.bo.driverName + "stats.csv")
	if err != nil {
		log.Fatal("Unable to create statfile", err)
		os.Exit(1)
	}
	defer file.Close()
	defer statfile.Close()

	writeIntro(file, data)
	if myPresets.bo.calcAverageActiveTime {
		_, err = file.WriteString("Average time of active read: " + fmt.Sprintf("%f", data.AverageTime) + "\n")
		if err != nil {
			log.Fatal("Unable to write in file", err)
		}
	}
	writeIntro(statfile, data)

	for _, currentDuration := range data.Durations {
		_, err = file.WriteString(fmt.Sprint(currentDuration) + "\n")
		if err != nil {
			log.Fatal("Unable to write in file", err)
		}
	}
	for _, stat := range data.Stats {
		for _, s := range stat {
			_, err = statfile.WriteString(s.CPU + "," + s.Memory.Percent + "\n")
			if err != nil {
				log.Fatal("Cannot write into file", err)
			}
		}
	}
}

func calcAverageTime(durations []int64) float64 {
	var at int64
	at = 0
	for _, d := range durations {
		at += d
	}
	return float64(at) / math.Max(float64(len(durations)), 1)
}
