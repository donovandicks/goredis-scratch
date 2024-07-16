package main

import (
	"flag"
	"log"

	"github.com/donovandicks/goredis/persistence"
	"github.com/donovandicks/goredis/server"
)

var (
	persistStrategy string
	persistFile     string
)

func main() {
	var strategy persistence.Strategy
	var err error
	flag.StringVar(&persistStrategy, "s", "aof", "the persistence strategy. defaults to 'aof' for append-only file")
	flag.StringVar(&persistFile, "f", "", "the path for an append-only file")

	switch persistStrategy {
	case "aof":
		var path string
		if persistFile == "" {
			path = "./aof"
		} else {
			path = persistFile
		}

		strategy, err = persistence.NewFile(path)
		if err != nil {
			log.Fatalf("failed to start AoF persistence strategy: %v", err)
		}
	default:
		log.Fatalf("invalid arguments")
	}

	server.NewServer(strategy).Run()
}
