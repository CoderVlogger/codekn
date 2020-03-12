package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"profx"
	"profx/storage/crawler"
	"strconv"
	"sync"
	"syscall"
	"time"
)

/*
	v0.3.*		crawler, periodic run (static sourcer)
	v0.4.0		dynamic sourcer and source rules
	v0.4.1		fix include rule match
	v0.5.0		cross ref info
	v0.5.1		fix nil pointer for db link update
	v0.5.2		fix to continue crawling if there is no source link to update
*/
var version = "v0.5.2"

var (
	fPing    = flag.Bool("ping", false, "db ping mode")
	fMigrate = flag.Bool("migrate", false, "run db migrations")
	fCrawl   = flag.Bool("crawl", false, "run crawler execution")

	wg sync.WaitGroup
)

func main() {
	// print version and parse flags
	log.Println(version)
	flag.Parse()

	config := crawler.Config{
		Username:  getEnv("DBC_USERNAME", "profx"),
		Password:  getEnv("DBC_PASSWORD", "profx"),
		ReadHost:  getEnv("DBC_READ_HOST", "localhost"),
		ReadPort:  getEnv("DBC_READ_PORT", "3306"),
		WriteHost: getEnv("DBC_WRITE_HOST", "localhost"),
		WritePort: getEnv("DBC_WRITE_PORT", "3306"),
		Schema:    getEnv("DBC_SCHEMA", "profx"),

		ReadMaxConn:   getEnv("DBC_READ_MAX_CONN", "50"),
		ReadIdleConn:  getEnv("DBC_READ_IDLE_CONN", "10"),
		WriteMaxConn:  getEnv("DBC_WRITE_MAX_CONN", "50"),
		WriteIdleConn: getEnv("DBC_WRITE_IDLE_CONN", "10"),
	}

	db, err := pingDB(config)
	if err != nil {
		log.Fatalf("database ping error: %v", err)
	} else {
		log.Println("database ping success")
	}

	if *fPing {
		log.Println("ping mode, bye!")
		return
	}

	if *fMigrate {
		err := crawler.Migrate(config)
		if err != nil {
			log.Fatalf("migration error")
		} else {
			log.Println("no migration errors")
			return
		}
	} else {
		log.Println("no flag provided, skipping migration")
	}

	if *fCrawl {
		log.Println("Starting crawler...")
		scp := profx.NewCollyScraper()
		src := profx.NewPersistentMemorySourcer(db)
		crw := profx.NewWebCrawler(src, scp, db)
		wg.Add(1)
		go crawl(&crw)
		log.Println("Crawler started")
		wg.Wait()
	} else {
		log.Println("no flag provided, skipping crawler")
	}
}

func getEnv(n, d string) string {
	v := os.Getenv(n)
	if v == "" {
		return d
	}
	return v
}

func pingDB(cfg crawler.Config) (profx.CrawlerRepository, error) {
	db, err := cfg.New()
	if err != nil {
		return nil, fmt.Errorf("can't check database %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database %v", err)
	} else {
		return db, nil
	}
}

func crawl(crw *profx.WebCrawler) {
	runCount := 10
	runInterval := 1 // seconds
	runCounter := 0

	rc, foundRC := syscall.Getenv("PROFX_RUN_COUNT")
	if foundRC {
		convRC, err := strconv.Atoi(rc)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if convRC >= 0 {
			runCount = convRC
		} else {
			log.Println("env variable PROFX_RUN_COUNT ignored, must be greater than -1")
		}
	}

	// TODO: Refactor to a func (issue #3).
	ri, foundRI := syscall.Getenv("PROFX_RUN_INTERVAL")
	if foundRI {
		convRI, err := strconv.Atoi(ri)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if convRI > 0 {
			runInterval = convRI
		} else {
			log.Println("env variable PROFX_RUN_INTERVAL ignored, must be greater than 0")
		}
	}

	running := false

	done := make(chan bool)
	run := make(chan bool)
	go func() {
		for {
			if runCount != 0 && runCounter >= runCount {
				done <- true
			}

			if !running {
				run <- true
			}
		}
	}()

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-done:
			log.Println("\r- run count reached")

			wg.Done()
			os.Exit(0)
		case <-interrupt:
			log.Println("\r- program interrupted")
			os.Exit(0)
		case <-run:
			running = true
			log.Println("Running Crawl()")

			err := crw.Crawl()
			// err := errors.New("asd")

			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Duration(runInterval) * time.Second)
			runCounter++
			running = false
		}
	}
}
