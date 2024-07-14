package main

import (
	"github.com/donovandicks/goredis/server"
)

func main() {
	server.NewServer().Run()
}
