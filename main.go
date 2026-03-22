package main

import (
	"flag"
	"log"

	// IMPORT YOUR OWN PACKAGES HERE! Replace YOUR_USERNAME.
	"github.com/Darsh0531/build-redis-go/config"
	"github.com/Darsh0531/build-redis-go/server"
)

func setupFlags() {
	// This reads command line arguments. If someone runs your code with `--port=8080`,
	// it will update the config.Port variable.
	// The "&" means we are passing the "memory address" of the variable so the flag package can modify it.
	flag.IntVar(&config.Port, "port", 7379, "port for the redis server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("rolling the dice ")
	server.RunSyncTCPServer() // We will build this function next!
}
