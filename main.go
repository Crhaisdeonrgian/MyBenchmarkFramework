package main

import (
	"database/sql"
	"log"
	"os"
)

func launchDataWorkers(data Data) {
	if myPresets.bo.makeFiles {
		writeData(data)
	}
	if myPresets.bo.drawPlots {
		writePlots(data)
	}
}

func launchBenchmarks(dbStd *sql.DB) {
	if myPresets.bo.runReadWithReadBackGround {
		data := ActiveRequestWithBackgroundLoadBenchmark(dbStd, executeReadWithTimeout, executeRead, myPresets.bo.BackGroundReadQuery, myPresets.bo.ActiveReadQuery)
		launchDataWorkers(data)
	}
	if myPresets.bo.runWriteWithWriteBackGround {
		data := ActiveRequestWithBackgroundLoadBenchmark(dbStd, executeWriteWithTimeout, executeWrite, myPresets.bo.BackGroundWriteQuery, myPresets.bo.ActiveWriteQuery)
		launchDataWorkers(data)
	}
	if myPresets.bo.runRead {
		data := ActiveRequestBenchmark(dbStd, executeRead, myPresets.bo.ActiveReadQuery)
		launchDataWorkers(data)
	}
	if myPresets.bo.runWrite {
		data := ActiveRequestBenchmark(dbStd, executeWrite, myPresets.bo.ActiveWriteQuery)
		launchDataWorkers(data)
	}

}

func main() {
	pr := SetOptions()
	mysqlContainer := SetConfig(pr)
	dbStd := connectToDB(pr.bo)
	launchBenchmarks(dbStd)

	if err := dockerPool.Purge(mysqlContainer); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(0)
}
