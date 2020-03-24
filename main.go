package main

import (
	"fmt"
	"log"
	"realtimechat/server"
)

func main() {

	serv := server.New()
	fmt.Println("running ...")
	if err := serv.Start(); err != nil {
		log.Fatalf("start error: %v", err)
	}

}