package main

import (
	"flash"
	"fmt"
	"os"
)

/*
	v0.3	env in json
	v0.4	add flash/app
	v0.5	add cache
	v0.6	gzip encoding
	v0.7	optimized gzip compression
*/
var version = "v0.6"

func main() {
	host := os.Getenv("FLASH_HOST")
	port := os.Getenv("FLASH_PORT")

	// start server
	server := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting server at %s\n", server)

	app := flash.NewApp(server, version)
	app.Run()
}
