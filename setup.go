package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type Presets struct {
	dbo dbOptions
	eo  envRunOptions
	bo  benchmarkOptions
}

type dbOptions struct {
	DBName   string
	RowCount int
	User     string
	Password string
}
type envRunOptions struct {
	FilePath      string
	MountPoints   []string
	CPUCount      int64
	Memory        int64
	Repository    string
	Tag           string
	Env           []string
	MaxDockerWait time.Duration
}
type benchmarkOptions struct {
	driverName                  string
	runReadWithReadBackGround   bool
	runWriteWithWriteBackGround bool
	runRead                     bool
	runWrite                    bool
	BackGroundReadQuery         string
	ActiveReadQuery             string
	BackGroundWriteQuery        string
	ActiveWriteQuery            string
	BackgroundContextTimeout    time.Duration
	BackgroundPeriod            time.Duration
	ActivePeriod                time.Duration
	CaptureDataPeriod           time.Duration
	BenchmarkTime               time.Duration
	drawPlots                   bool
	makeFiles                   bool
	calcAverageActiveTime       bool
}

// nolint:gochecknoglobals
var dockerPool *dockertest.Pool // the connection to docker
// nolint:gochecknoglobals
var systemdb *sql.DB // the connection to the mysql 'system' database
// nolint:gochecknoglobals
var sqlConfig *mysql.Config // the mysql container and config for connecting to other databases
// nolint:gochecknoglobals
var testMu *sync.Mutex // controls access to sqlConfig

var myPresets Presets

func SetOptions() Presets {
	//Options of environment to run in:
	eo := envRunOptions{
		FilePath:      "/Users/igorvozhga/DIPLOMA/",                                  //where to save results
		MountPoints:   []string{"/Users/igorvozhga/DIPLOMA/mountDir:/var/lib/mysql"}, //where to mount persistent volume
		CPUCount:      1,                                                             //bound CPU to load volume
		Memory:        1024 * 1024 * 1024 * 1,                                        //1Gb bound memory usage to load volume
		Repository:    "mysql",
		Tag:           "5.6",
		Env:           []string{"MYSQL_ROOT_PASSWORD=secret"},
		MaxDockerWait: time.Minute * 2, //how long wait for connection
	}
	//Options of testing db:
	dbo := dbOptions{
		User:     "root",
		Password: "secret",
		DBName:   "MyDB", //It's on your own
		RowCount: 100000, //Size of table(we're looking only on one table for now)
	}
	/*
		Here you can choose what benchmark to run,
		what parameters to use
		and what to do with the results
	*/
	bo := benchmarkOptions{
		driverName:                "mysql",
		runReadWithReadBackGround: true,
		BackGroundReadQuery:       SlowRead,
		ActiveReadQuery:           MediumRead,
		BackgroundContextTimeout:  15 * time.Second,
		BackgroundPeriod:          5 * time.Second,
		ActivePeriod:              2 * time.Second,
		CaptureDataPeriod:         1 * time.Second,
		BenchmarkTime:             180 * time.Second,
		drawPlots:                 true,
		makeFiles:                 true,
		calcAverageActiveTime:     true,
	}
	myPresets = Presets{
		bo:  bo,
		eo:  eo,
		dbo: dbo,
	}
	return myPresets
}

func SetConfig(pr Presets) *dockertest.Resource {
	_ = mysql.SetLogger(log.New(ioutil.Discard, "", 0)) // silence mysql logger
	testMu = &sync.Mutex{}

	var err error
	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	dockerPool.MaxWait = time.Minute * 2

	runOptions := dockertest.RunOptions{
		Repository: pr.eo.Repository,
		Tag:        pr.eo.Tag,
		Env:        pr.eo.Env,
		Mounts:     pr.eo.MountPoints,
	}
	mysqlContainer, err := dockerPool.RunWithOptions(&runOptions, func(hostcfg *docker.HostConfig) {
		hostcfg.CPUCount = pr.eo.CPUCount
		//hostcfg.CPUPercent = 100
		hostcfg.Memory = pr.eo.Memory
	})
	if err != nil {
		log.Fatalf("could not start mysqlContainer: %s", err)
	}
	sqlConfig = &mysql.Config{
		User:                 pr.dbo.User,
		Passwd:               pr.dbo.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("localhost:%s", mysqlContainer.GetPort("3306/tcp")),
		DBName:               pr.dbo.DBName,
		AllowNativePasswords: true,
	}

	if err = dockerPool.Retry(func() error {
		systemdb, err = sql.Open(pr.bo.driverName, sqlConfig.FormatDSN())
		if err != nil {
			return err
		}
		return systemdb.Ping()
	}); err != nil {
		log.Fatal(err)
	}

	return mysqlContainer
}
