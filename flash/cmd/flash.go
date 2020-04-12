package main

import (
	"flash/internal/app"
	"flash/internal/infra/storage/crawler"
	"fmt"
	"log"
	"os"
)

/*
	v0.3	env in json
	v0.4	add flash/app
	v0.5	add cache
	v0.6	gzip encoding
	v0.7	optimized gzip compression
	v0.8	index page and /api route
*/
var version = "v0.8"

func main() {
	host := os.Getenv("FLASH_HOST")
	port := os.Getenv("FLASH_PORT")

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

	db, err := config.New()
	if err != nil {
		log.Fatalf("database error: %v", err)
	} else {
		log.Println("database success")
	}

	// start server
	server := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting server at %s\n", server)

	appIns := app.NewApp(server, version, db)
	appIns.Run()
}

func getEnv(n, d string) string {
	v := os.Getenv(n)
	if v == "" {
		return d
	}
	return v
}
