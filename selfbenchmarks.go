package main

import (
	"database/sql"
	"github.com/KyleBanks/dockerstats"
	"log"
	"time"
)

type BenchmarkModel struct {
	BenchmarkType   string
	ActiveQuery     string
	BackgroundQuery string
}

type Data struct {
	Benchmark   BenchmarkModel
	AverageTime float64
	Durations   []int64
	Stats       [][]dockerstats.Stats
	DurationX   []float64
	DurationY   []float64
}

func organizeData(durations chan int64, stats chan []dockerstats.Stats) Data {
	var data Data
	var durationV []int64
	var statsV [][]dockerstats.Stats
	for currentDuration := range durations {
		durationV = append(durationV, currentDuration)
	}
	for stat := range stats {
		statsV = append(statsV, stat)
	}
	if myPresets.bo.calcAverageActiveTime {
		data.AverageTime = calcAverageTime(durationV)
	}
	data.Durations = durationV
	data.Stats = statsV
	return data
}

//TODO: Реализовать real-life бенч и поискать еще паттерны

func ActiveRequestWithBackgroundLoadBenchmark(dbStd *sql.DB, backgroundFunc func(*sql.DB, string, time.Duration), activeFunc func(*sql.DB, string), backgroundQuery string, activeQuery string) Data {
	var data Data
	data.Benchmark = BenchmarkModel{
		BenchmarkType:   "ActiveRequestWithBackgroundLoad",
		ActiveQuery:     myPresets.bo.ActiveReadQuery,
		BackgroundQuery: myPresets.bo.BackGroundReadQuery,
	}
	var durations = make(chan int64, int(myPresets.bo.BenchmarkTime)/int(myPresets.bo.ActivePeriod))
	var done = make(chan struct{})
	var testEnded = make(chan struct{})
	var stats = make(chan []dockerstats.Stats, int(myPresets.bo.BenchmarkTime)/int(myPresets.bo.CaptureDataPeriod))
	backgroundTicker := time.NewTicker(myPresets.bo.BackgroundPeriod)
	activeTicker := time.NewTicker(myPresets.bo.ActivePeriod)
	statTicker := time.NewTicker(myPresets.bo.CaptureDataPeriod)

	go func(chan int64, chan struct{}, chan struct{}, chan []dockerstats.Stats) {
		for {
			select {
			case <-statTicker.C:
				s, err := dockerstats.Current()
				if err != nil {
					log.Println("Unable to get stats ", err)
				}
				stats <- s
			case <-done:
				close(durations)
				close(stats)
				close(testEnded)
				return
			case <-backgroundTicker.C:
				go func() {
					backgroundFunc(dbStd, backgroundQuery, myPresets.bo.BackgroundContextTimeout)
				}()
			case <-activeTicker.C:
				go func(chan int64, chan struct{}) {
					start := time.Now()
					activeFunc(dbStd, activeQuery)
					select {
					case _, is_open := <-testEnded:
						if is_open {
							// no one is putting anything
							log.Fatal("how u did it??")
						} else {
							// chan closed doing nothing
						}
					default:
						// chan is open so test is running
						d := time.Since(start).Milliseconds()
						log.Println("MediumQuery duration: ", d)
						durations <- d
					}
				}(durations, testEnded)
			}
		}
	}(durations, done, testEnded, stats)
	time.Sleep(myPresets.bo.BenchmarkTime)
	backgroundTicker.Stop()
	activeTicker.Stop()
	statTicker.Stop()
	done <- struct{}{}
	return organizeData(durations, stats)
}
func ActiveRequestBenchmark(dbStd *sql.DB, activeFunc func(*sql.DB, string), activeQuery string) Data {
	var durations = make(chan int64, int(myPresets.bo.BenchmarkTime)/int(myPresets.bo.ActivePeriod))
	var done = make(chan struct{})
	var testEnded = make(chan struct{})
	var stats = make(chan []dockerstats.Stats, int(myPresets.bo.BenchmarkTime)/int(myPresets.bo.CaptureDataPeriod))
	ticker := time.NewTicker(myPresets.bo.ActivePeriod)
	statTicker := time.NewTicker(myPresets.bo.CaptureDataPeriod)
	go func(chan int64, chan struct{}, chan struct{}, chan []dockerstats.Stats) {
		for {
			select {
			case <-statTicker.C:
				s, err := dockerstats.Current()
				if err != nil {
					log.Println("Unable to get stats ", err)
				}
				stats <- s
			case <-done:
				close(durations)
				close(stats)
				close(testEnded)
				return
			case <-ticker.C:
				go func(chan int64, chan struct{}) {
					start := time.Now()
					activeFunc(dbStd, activeQuery)
					select {
					case _, is_open := <-testEnded:
						if is_open {
							// no one is putting anything
							log.Fatal("how u did it??")
						} else {
							// chan closed doing nothing
						}
					default:
						// chan is open so test is running
						d := time.Since(start).Milliseconds()
						log.Println("MediumQuery duration: ", d)
						durations <- d
					}
				}(durations, testEnded)
			}
		}
	}(durations, done, testEnded, stats)
	time.Sleep(myPresets.bo.BenchmarkTime)
	ticker.Stop()
	statTicker.Stop()
	done <- struct{}{}
	return organizeData(durations, stats)
}
